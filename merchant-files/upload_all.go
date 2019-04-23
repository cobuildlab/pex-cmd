package merchants

import (
	"fmt"
	"log"
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
			log.Println("The process could not be executed because of the existence of the file upload.lock")
			os.Exit(0)
		}

		emptyFile, err := os.Create("./upload.lock")
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		emptyFile.Close()

		fileList, err := UploadList()
		if err != nil {
			log.Println(err)
			return
		}

		for i := 0; i < len(fileList); i++ {
			v := fileList[i]

			log.Println("╭─Uploading", v.Name, "...", fmt.Sprintf("%d/%d", i+1, len(fileList)))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v.Name, Verbose)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println("├─⇢ ...", "Uploaded!", v.Name)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()

		}

		err = os.Remove("./upload.lock")
		if err != nil {
			return
		}
	},
}
