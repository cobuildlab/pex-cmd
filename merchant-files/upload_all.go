package merchants

import (
	"fmt"
	"os"
	"time"

	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdUploadAll Command to upload all the merchant files to the database<Paste>
var CmdUploadAll = &cobra.Command{
	Use:   "all",
	Short: "Upload all the merchants files to the database",
	Long:  "Upload all the merchants files to the database",
	Run: func(cmd *cobra.Command, args []string) {
		exist, err := utils.CheckExistence("./upload.lock")
		if err != nil {
			return
		}

		if exist {
			fmt.Println(TimeNow(), "The process could not be executed because of the existence of the file upload.lock")
			os.Exit(0)
		}

		fileList, err := UploadList()
		if err != nil {
			fmt.Println(TimeNow(), err)
			return
		}

		for i, v := range fileList {
			fmt.Println(TimeNow(), "╭─Uploading", v.Name, "...", fmt.Sprintf("%d/%d", i+1, len(fileList)))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v.Name, Verbose)
			if err != nil {
				fmt.Println(TimeNow(), err)
				return
			}
			fmt.Println(TimeNow(), "├─⇢ ...", "Uploaded!", v.Name)
			fmt.Println(TimeNow(), "│")
			fmt.Println(TimeNow(), "├──⇢ Uploaded products:", totalProductsUpload)
			fmt.Println(TimeNow(), "├──⇢ Updated products:", totalProductsUpdated)
			fmt.Println(TimeNow(), "╰──⇢ Duration:", time.Since(start))
			fmt.Println()

			err = os.Remove("./upload.lock")
			if err != nil {
				return
			}
		}
	},
}
