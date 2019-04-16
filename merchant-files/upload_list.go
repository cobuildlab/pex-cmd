package merchants

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdUploadList Command to list the available merchants files to upload
var CmdUploadList = &cobra.Command{
	Use:   "list",
	Short: "List of merchants files available to upload to the database",
	Long:  "List of merchants files available to upload to the database",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, err := UploadList()
		if err != nil {
			fmt.Println(err)
			return
		}

		var count uint64
		var totalSize int64
		for _, v := range fileList {
			count++
			totalSize += v.Size
			fmt.Println("╭─Filename:", v.Name)
			fmt.Println("├─⇢ Size:", v.Size)
			fmt.Println("╰─⇢ Time:", v.ModTime)
			fmt.Println()
		}
		fmt.Println("╭─Total files:", count)
		fmt.Println("╰─Total size:", totalSize)
	},
}

//UploadList List the available merchants files to upload
func UploadList() (filesXML []utils.FileInfo, err error) {
	files, err := ioutil.ReadDir(DecompressPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == MerchantFileFormatExt {
				var fileXML utils.FileInfo
				fileXML.Name = file.Name()
				fileXML.Size = file.Size()
				fileXML.ModTime = file.ModTime()

				filesXML = append(filesXML, fileXML)
			}
		}
	}

	return
}
