package merchants

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cobuildlab/pex-cmd/errors"
	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdDownloadAll Command to download all the merchant files
var CmdDownloadAll = &cobra.Command{
	Use:   "all",
	Short: "Download all the merchants files of the FTP server",
	Long:  "Download all the merchants files of the FTP server",
	Run: func(cmd *cobra.Command, args []string) {
		exist, err := utils.CheckExistence("./download.lock")
		if err != nil {
			return
		}

		if exist {
			fmt.Println(TimeNow(), "The process could not be executed because of the existence of the file download.lock")
			os.Exit(0)
		}

		fmt.Println(TimeNow(), "* Limit file size:", LimitSize)

		err = DownloadAll(LimitSize)
		if err != nil {
			fmt.Println(TimeNow(), err)
			return
		}
	},
}

//DownloadAll Download all merchant files
func DownloadAll(limitSize uint64) (err error) {
	var countDownloadFiles, countDecompressFiles, countFailedDownloads uint64

	if ok, _ := utils.CheckExistence(utils.FTPPathFiles); !ok {
		err = os.MkdirAll(utils.FTPPathFiles, 0777)
		if err != nil {
			err = errors.ErrorPathNotExist

			return
		}
	}

	client, err := utils.GetConnectionFTP(
		utils.FTPHost, utils.FTPPort,
		utils.FTPUsername, utils.FTPPassword,
	)

	if err != nil {
		return
	}

	defer client.Quit()
	defer client.Logout()

	entries, err := client.List(FTPWildcardFilesFilter)
	if err != nil {
		err = errors.ErrorConsultingData

		return
	}

	start := time.Now()
	for i, entry := range entries {
		if entry.Type == 0 {
			if entry.Size > limitSize {
				fmt.Println(TimeNow(), "[x]", entry.Name, "Out of size limit:", entry.Size, fmt.Sprintf("%d/%d", i+1, len(entries)))
				fmt.Println()
			} else {
				fileDownload := filepath.Join(utils.FTPPathFiles, entry.Name)

				fileExt := filepath.Ext(entry.Name)

				if fileExt == MerchantFileCompressExt {

					fmt.Println(TimeNow(), "╭─Downloading", entry.Name, "...", fmt.Sprintf("%d/%d", i+1, len(entries)))
					fmt.Println(TimeNow(), "├─⇢ Size:", entry.Size)
					if err = utils.DownloadGzipFileFTP(entry.Name, utils.FTPPathFiles); err != nil {
						fmt.Println(TimeNow(), "[x]", "An error has occurred downloading the file:", entry.Name, fmt.Sprintf("%d/%d", i+1, len(entries)))
						fmt.Println()
						countFailedDownloads++
						continue
					}
					countDownloadFiles++

					if err = utils.DecompressFileGzip(fileDownload, DecompressPath); err != nil {
						return
					}
					countDecompressFiles++

					fmt.Println(TimeNow(), "╰─...", "Discharged", entry.Name)
					fmt.Println()
				}
			}

		}
	}

	fmt.Println(TimeNow(), "├─⇢ ...", "Discharged!")
	fmt.Println(TimeNow(), "│")
	fmt.Println(TimeNow(), "├──⇢ Downloaded files:", countDownloadFiles)
	fmt.Println(TimeNow(), "├──⇢ Failed downloads:", countFailedDownloads)
	fmt.Println(TimeNow(), "├──⇢ Unzipped files:", countDownloadFiles)
	fmt.Println(TimeNow(), "╰──⇢ Duration:", time.Since(start))

	err = os.Remove("./download.lock")
	if err != nil {
		return
	}

	return
}
