package merchants

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cobuildlab/pex-cmd/databases"
	"github.com/cobuildlab/pex-cmd/models"
	"github.com/cobuildlab/pex-cmd/services/products"
	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

//CmdUploadFile Command to upload a merchant file to the database
var CmdUploadFile = &cobra.Command{
	Use:   "file [string]",
	Short: "Upload a merchant file to the database",
	Long:  "Upload a merchant file to the database",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i, v := range args {
			log.Println("╭─Uploading", v, "...", fmt.Sprintf("%d/%d", i+1, len(args)))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, totalProductsFailed, err := UploadFile(v, Verbose)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("├─⇢ ...", "Uploaded!", v)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("├──⇢ Failed products:", totalProductsFailed)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()
		}
	},
}

//UploadFile Upload a merchant file to the database
func UploadFile(filename string, verbose bool) (totalProductsUpload, totalProductsUpdated, totalProductsFailed uint64, err error) {
	var merchant utils.Merchant
	var product models.Product
	var wg sync.WaitGroup

	pathFile := filepath.Join(DecompressPath, filename)

	clientDB, err := databases.NewClient(databases.Username, databases.Password)
	if err != nil {
		return
	}

	fileXML, err := os.Open(pathFile)
	if err != nil {
		return
	}

	log.Println("├─⇢ Product Counter:", CountProductsInMerchantFile(fileXML))
	fileXML.Close()

	dbMerchants := databases.OpenDB(clientDB, databases.DBNameMerchants)
	dbProducts := databases.OpenDB(clientDB, databases.DBNameProducts)

	productsRemoteID, err := products.GetIDs()
	if err != nil {
		return
	}

	fileXML, err = os.Open(pathFile)
	if err != nil {
		return
	}
	merchant, err = UploadMerchantInMerchantFile(fileXML, dbMerchants)
	if err != nil {
		return
	}
	fileXML.Close()

	fileXML, err = os.Open(pathFile)
	if err != nil {
		return
	}
	defer fileXML.Close()

	dec := xml.NewDecoder(fileXML)

	for {
		token, err := dec.Token()
		if err == io.EOF {
			err = nil
			break
		}

		switch token.(type) {
		case xml.StartElement:
			start := token.(xml.StartElement)

			switch start.Name.Local {
			case "product":
				//Save products

				databases.QueueWriter.Wait()

				if err = dec.DecodeElement(&product, &start); err != nil {
					break
				}

				if product.Price.Currency != "" && product.Price.Currency != "USD" {
					break
				}

				if product.Discount.Currency != "" && product.Discount.Currency != "USD" {
					break
				}

				wg.Add(1)
				go uploadProduct(product, merchant, dbProducts, &totalProductsUpload, &totalProductsUpdated, &totalProductsFailed, productsRemoteID, &wg)
			}
		}
	}

	wg.Wait()

	os.Remove("data/rakuten/decompress/" + filename)
	os.Remove("data/rakuten/" + filename + ".gz")

	return
}

func uploadProduct(product models.Product, merchantLocal utils.Merchant, dbProducts databases.DB, totalProductsUpload, totalProductsUpdated, totalProductsFailed *uint64, productsRemoteID []string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	var countFailed uint
	var exists bool
	var productRemote models.Product
	var productRemoteREV string

	product.Merchant = models.Merchant{
		ID:   merchantLocal.MerchantID,
		Name: merchantLocal.MerchantName,
	}

ForProduct:
	for {

		for _, v := range productsRemoteID {
			if v == product.ID {

				var result interface{}
				err = databases.ReadElement(dbProducts, product.ID, &result, models.OptionsDB{})
				if err != nil {
					if countFailed >= 5 {
						continue ForProduct
					}
					atomic.AddUint64(totalProductsFailed, 1)
					return

				}
				if result != nil {
					exists = true

					productRemoteREV = result.(map[string]interface{})["_rev"].(string)
					if err = mapstructure.Decode(result, &productRemote); err != nil {
						if countFailed >= 5 {
							continue ForProduct
						}
						return
					}
					productRemote.ID = result.(map[string]interface{})["_id"].(string)

				}
				break
			}
		}

		if err != nil {
			if Verbose {
				log.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), product.ID)
			}
			if countFailed >= 5 {
				countFailed++
				continue ForProduct
			}
			atomic.AddUint64(totalProductsFailed, 1)
			return
		}

		if exists {
			if productRemote != product {
				_, err = databases.UpdateElement(dbProducts, product.ID, productRemoteREV, product)
				if err != nil {
					if Verbose {
						log.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), product.ID)
					}
					if countFailed >= 5 {
						countFailed++
						continue ForProduct
					}
					atomic.AddUint64(totalProductsFailed, 1)
					return
				}
				if Verbose {
					log.Printf("├──⇢+ Success: Product #%s, Updated\n", product.ID)
				}
				atomic.AddUint64(totalProductsUpdated, 1)
			}
		} else {
			_, _, err = databases.CreateElement(dbProducts, product)
			if err != nil {
				if Verbose {
					log.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), product.ID)
				}
				if countFailed >= 5 {
					countFailed++
					continue ForProduct
				}
				atomic.AddUint64(totalProductsFailed, 1)
				return
			}
			if Verbose {
				log.Printf("├──⇢+ Success: Product #%s, Uploaded\n", product.ID)
			}
			atomic.AddUint64(totalProductsUpload, 1)
		}
		break
	}

	return

}
