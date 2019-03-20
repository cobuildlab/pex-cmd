package merchants

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/4geeks/pex-cmd/errors"
	"github.com/4geeks/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdDownloadFile Command to download a merchant file
var CmdDownloadFile = &cobra.Command{
	Use:   "file [string]",
	Short: "Download a merchant file from the FTP server",
	Long:  "Download a merchant file from the FTP server",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, v := range args {
			fmt.Println("╭─Downloading", v, "...", fmt.Sprintf("%d/%d", i+1, len(args)))
			err := DownloadFile(v)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("╰─...", "Discharged", v)
			fmt.Println()
		}
	},
}

//DownloadFile Download a merchant file
func DownloadFile(filename string) (err error) {
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

	fileDownload := filepath.Join(utils.FTPPathFiles, filename)

	if err = utils.DownloadGzipFileFTP(client, filename, utils.FTPPathFiles); err != nil {
		return
	}

	err = utils.DecompressFileGzip(fileDownload, DecompressPath)

	return
}
