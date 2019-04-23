package merchants

import (
	"log"

	"github.com/cobuildlab/pex-cmd/errors"
	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdDownloadList Command to list merchants files on the FTP server
var CmdDownloadList = &cobra.Command{
	Use:   "list",
	Short: "List of merchants files on FTP server",
	Long:  "List of merchants files on FTP server",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, err := DownloadList()
		if err != nil {
			log.Println(err)
			return
		}
		var count uint64
		var totalSize int64
		for _, v := range fileList {
			count++
			totalSize += v.Size
			log.Println("╭─Filename:", v.Name)
			log.Println("├─⇢ Size:", v.Size)
			log.Println("╰─⇢ Time:", v.ModTime)
			log.Println()
		}
		log.Println("╭─Total files:", count)
		log.Println("╰─Total size:", totalSize)
		log.Println("*The total size is the sum of the size of all files on the FTP server and should not be used as an exact reference*")
	},
}

//DownloadList List merchants files on the FTP server
func DownloadList() (entryList []utils.FileInfo, err error) {
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

	for _, v := range entries {
		entry := utils.FileInfo{
			Name:    v.Name,
			Size:    int64(v.Size),
			ModTime: v.Time,
		}

		entryList = append(entryList, entry)
	}

	return
}
