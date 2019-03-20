package merchants

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

//CmdUploadAll Command to upload all the merchant files to the database<Paste>
var CmdUploadAll = &cobra.Command{
	Use:   "all",
	Short: "Upload all the merchants files to the database",
	Long:  "Upload all the merchants files to the database",
	Run: func(cmd *cobra.Command, args []string) {
		fileList, err := UploadList()
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, v := range fileList {
			fmt.Println("╭─Uploading", v.Name, "...", fmt.Sprintf("%d/%d", i+1, len(fileList)))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v.Name, Verbose)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("├─⇢ ...", "Uploaded!", v.Name)
			fmt.Println("│")
			fmt.Println("├──⇢ Uploaded products:", totalProductsUpload)
			fmt.Println("├──⇢ Updated products:", totalProductsUpdated)
			fmt.Println("╰──⇢ Duration:", time.Since(start))
			fmt.Println()
		}
	},
}
