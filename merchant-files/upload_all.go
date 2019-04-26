package merchants

import (
	"fmt"
	"log"
	"os"
	"sort"
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
		f, err := os.OpenFile("upload.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.SetOutput(f)

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

		sort.Slice(fileList, func(i, j int) bool {
			return fileList[i].Name > fileList[j].Name
		})

		var countFiles int = len(fileList)
		var i int
		for len(fileList) != 0 {
			if len(fileList) != countFiles {
				countFiles += countFiles - len(fileList)
			}

			v := fileList[0]

			log.Println("╭─Uploading", v.Name, "...", fmt.Sprintf("%d/%d", i+1, countFiles))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v.Name, Verbose)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("├─⇢ ...", "Uploaded!", v.Name)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()

			fileList, err = UploadList()
			if err != nil {
				log.Println(err)
				return
			}

			i++
		}

		err = os.Remove("./upload.lock")
		if err != nil {
			return
		}
	},
}
