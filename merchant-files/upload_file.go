package merchants

import (
	"encoding/xml"
	"fmt"
	"github.com/cobuildlab/pex-cmd/services/products"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cobuildlab/pex-cmd/databases"
	"github.com/cobuildlab/pex-cmd/models"
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

		productsRemoteID, err := products.GetIDs()
		if err != nil {
			log.Println("ERROR Obtaining the Product Ids", err)
			return
		}

		for i, v := range args {
			log.Println("╭─Uploading", v, "...", fmt.Sprintf("%d/%d", i+1, len(args)))

			start := time.Now()
			totalProductsUpload, totalProductsUpdated, totalProductsFailed, err := UploadFile(v, productsRemoteID, Verbose)
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
func UploadFile(filename string, productsRemoteID []string, verbose bool) (totalProductsUpload, totalProductsUpdated, totalProductsFailed uint64, err error) {
	var merchant utils.Merchant
	var product models.Product

	pathFile := filepath.Join(DecompressPath, filename)

	clientDB, err := databases.NewClient(databases.Username, databases.Password)
	if err != nil {
		return
	}

	dbMerchants := databases.OpenDB(clientDB, databases.DBNameMerchants)
	dbProducts := databases.OpenDB(clientDB, databases.DBNameProducts)

	log.Println("├─⇢ Opening:", pathFile)
	fileXML, err := os.Open(pathFile)
	if err != nil {
		log.Println("├─⇢ Error Opening File:")
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

	// Loading Products
	dec := xml.NewDecoder(fileXML)
	var wg sync.WaitGroup
	productsDebuggerMap := make(map[string]int64)
	routinesCounter := 0
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
				go uploadProduct(product, merchant, dbProducts, &totalProductsUpload, &totalProductsUpdated, &totalProductsFailed, productsRemoteID, &wg, pathFile, productsDebuggerMap)

				if routinesCounter%5 == 0 {
					wg.Wait()
					time.Sleep(1 * time.Second)
				}
			}
		}
	}

	wg.Wait()

	os.Remove("data/rakuten/decompress/" + filename)
	os.Remove("data/rakuten/" + filename + ".gz")
	return
}

// Upload a product to the database
func uploadProduct(product models.Product, merchantLocal utils.Merchant, dbProducts databases.DB, totalProductsUpload, totalProductsUpdated, totalProductsFailed *uint64, productsRemoteID []string, wg *sync.WaitGroup, pathFile string, productMap map[string]int64) (err error) {
	defer wg.Done()

	var exists bool
	var productRemote models.Product
	var productRemoteREV string

	product.Merchant = models.Merchant{
		ID:   merchantLocal.MerchantID,
		Name: merchantLocal.MerchantName,
	}

	for _, v := range productsRemoteID {
		if v == product.ID {

			var result interface{}
			err = databases.ReadElement(dbProducts, product.ID, &result, models.OptionsDB{})
			if err != nil {
				log.Printf("├──⇢+ Error: %s Product #%s", err.Error(), product.ID)
				return
			}

			if result != nil {
				exists = true
				productRemoteREV = result.(map[string]interface{})["_rev"].(string)
				if err = mapstructure.Decode(result, &productRemote); err != nil {
					log.Printf("├──⇢+ Error on productRemoteREV: %s Product #%s", err.Error(), product.ID)
					return
				}
				productRemote.ID = result.(map[string]interface{})["_id"].(string)
			}
			break
		}
	}

	if exists {
		if productRemote != product {
			_, err = databases.UpdateElement(dbProducts, product.ID, productRemoteREV, product)
			if err != nil {
				log.Printf("├──⇢+ Error on UpdateElement: %s Product #%s, Retrying\n", err.Error(), product.ID)
				return
			}
			log.Printf("├──⇢+ Success: Product Updated #%s, Updated from %s -- opertaions on this product ID : %s \n", product.ID, pathFile, productMap[product.ID])
		}
		atomic.AddUint64(totalProductsUpdated, 1)
	} else {
		_, _, err = databases.CreateElement(dbProducts, product)
		if err != nil {
			log.Printf("├──⇢+ Error on CreateElement: %s Product #%s ", err.Error(), product.ID)
			return
		}
		log.Printf("├──⇢+ Success: Product Created #%s, Uploaded from %s -- opertaions on this product ID : %s \n", product.ID, pathFile, productMap[product.ID])
		atomic.AddUint64(totalProductsUpload, 1)
	}
	incrementProductsMap(productMap, product.ID)
	return
}

func incrementProductsMap(productsMap map[string]int64, productId string) {
	val, exist := productsMap[productId]
	if exist == true {
		val++
		productsMap[productId] = val
	} else {
		productsMap[productId] = 1
	}
}