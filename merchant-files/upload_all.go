package merchants

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
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
			pid, _ := ioutil.ReadFile("./upload.lock")
			log.Println("The process could not be executed because of the existence of the file upload.lock, PID From the creator : [" + string(pid) + "]")
			os.Exit(0)
		}

		lockFile, err := os.Create("./upload.lock")
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		lockFile.WriteString(strconv.Itoa(os.Getpid()))
		lockFile.Close()

		fileList, err := UploadList()
		if err != nil {
			log.Println(err)
			return
		}

		sort.Slice(fileList, func(i, j int) bool {
			return fileList[i].Name > fileList[j].Name
		})

		for i := 0; len(fileList) != 0; i++ {
			var countFiles int = len(fileList)

			v := fileList[0]

			log.Println("╭─Uploading", v.Name, "...", fmt.Sprintf("%d/%d", i+1, countFiles))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, totalProductsFailed, err := UploadFile(v.Name, Verbose)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("├─⇢ ...", "Uploaded!", v.Name)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("├──⇢ Failed products:", totalProductsFailed)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()

			fileList, err = UploadList()
			if err != nil {
				log.Println(err)
				return
			}
		}

		err = os.Remove("./upload.lock")
		if err != nil {
			return
		}
	},
}
