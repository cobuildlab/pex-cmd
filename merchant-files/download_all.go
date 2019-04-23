package merchants

import (
	"fmt"
	"log"
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
			log.Println("The process could not be executed because of the existence of the file download.lock")
			os.Exit(0)
		}

		emptyFile, err := os.Create("./download.lock")
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		emptyFile.Close()

		log.Println("* Limit file size:", LimitSize)

		err = DownloadAll(LimitSize)
		if err != nil {
			log.Println(err)
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
				log.Println("[x]", entry.Name, "Out of size limit:", entry.Size, fmt.Sprintf("%d/%d", i+1, len(entries)))
				log.Println()
			} else {
				fileDownload := filepath.Join(utils.FTPPathFiles, entry.Name)

				fileExt := filepath.Ext(entry.Name)

				if fileExt == MerchantFileCompressExt {

					log.Println("╭─Downloading", entry.Name, "...", fmt.Sprintf("%d/%d", i+1, len(entries)))
					log.Println("├─⇢ Size:", entry.Size)
					if err = utils.DownloadGzipFileFTP(entry.Name, utils.FTPPathFiles); err != nil {
						log.Println("[x]", "An error has occurred downloading the file:", entry.Name, fmt.Sprintf("%d/%d", i+1, len(entries)))
						log.Println()
						countFailedDownloads++
						continue
					}
					countDownloadFiles++

					if err = utils.DecompressFileGzip(fileDownload, DecompressPath); err != nil {
						return
					}
					countDecompressFiles++

					log.Println("╰─...", "Discharged", entry.Name)
					log.Println()
				}
			}

		}
	}

	log.Println("├─⇢ ...", "Discharged!")
	log.Println("│")
	log.Println("├──⇢ Downloaded files:", countDownloadFiles)
	log.Println("├──⇢ Failed downloads:", countFailedDownloads)
	log.Println("├──⇢ Unzipped files:", countDownloadFiles)
	log.Println("╰──⇢ Duration:", time.Since(start))

	err = os.Remove("./download.lock")
	if err != nil {
		return
	}

	return
}
