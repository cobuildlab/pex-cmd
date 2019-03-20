package merchants

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4geeks/pex-cmd/errors"
	"github.com/4geeks/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdDownloadAll Command to download all the merchant files
var CmdDownloadAll = &cobra.Command{
	Use:   "all",
	Short: "Download all the merchants files of the FTP server",
	Long:  "Download all the merchants files of the FTP server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("* Limit file size:", LimitSize)
		err := DownloadAll(LimitSize)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

//DownloadAll Download all merchant files
func DownloadAll(limitSize uint64) (err error) {
	var countDownloadFiles, countDecompressFiles uint64

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

	for i, entry := range entries {
		if entry.Type == 0 {
			if entry.Size > limitSize {
				fmt.Println("[x]", entry.Name, "Out of size limit:", entry.Size, fmt.Sprintf("%d/%d", i+1, len(entries)))
				fmt.Println()
			} else {
				fileDownload := filepath.Join(utils.FTPPathFiles, entry.Name)

				fileExt := filepath.Ext(entry.Name)

				if fileExt == MerchantFileCompressExt {

					fmt.Println("╭─Downloading", entry.Name, "...", fmt.Sprintf("%d/%d", i+1, len(entries)))
					fmt.Println("├─⇢ Size:", entry.Size)
					if err = utils.DownloadGzipFileFTP(client, entry.Name, utils.FTPPathFiles); err != nil {
						return
					}
					countDownloadFiles++

					if err = utils.DecompressFileGzip(fileDownload, DecompressPath); err != nil {
						return
					}
					countDecompressFiles++

					fmt.Println("╰─...", "Discharged", entry.Name)
					fmt.Println()
				}
			}

		}
	}

	return
}
