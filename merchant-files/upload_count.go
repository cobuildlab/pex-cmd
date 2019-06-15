package merchants

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

//CmdUploadCount Count the products of all the merchant files
var CmdUploadCount = &cobra.Command{
	Use:   "count",
	Short: "Count the products of all the merchant files",
	Long:  "Count the products of all the merchant files",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.OpenFile("count.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		log.SetOutput(f)

		fileList, err := UploadList()
		if err != nil {
			log.Println(err)
			return
		}

		var count int
		var size int64
		for _, v := range fileList {
			file, err := os.Open("data/rakuten/decompress/" + v.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
			countP, err := CountProductsInMerchantFile(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			log.Println("╭─Filename:", v.Name)
			log.Println("├─⇢ Size:", v.Size)
			log.Println("├─⇢ Time:", v.ModTime)
			log.Println("╰─⇢ Products:", countP)
			log.Println()

			count += countP
			size += v.Size
		}

		log.Println()
		log.Println("╭─ Total products:", count)
		log.Println("├─⇢ Total size:", size)
		log.Println("╰⇢ Total files:", len(fileList))
	},
}
