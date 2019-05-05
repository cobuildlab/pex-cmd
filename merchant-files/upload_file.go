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

var queueUpload = make(chan bool, 100000)
var ch chan int

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

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v, Verbose)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("├─⇢ ...", "Uploaded!", v)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()
		}
	},
}

//UploadFile Upload a merchant file to the database
func UploadFile(filename string, verbose bool) (totalProductsUpload, totalProductsUpdated uint64, err error) {
	pathFile := filepath.Join(DecompressPath, filename)

	clientDB, err := databases.NewClient(
		databases.Username,
		databases.Password,
	)

	if err != nil {
		return
	}

	var merchant utils.Merchant
	var merchantLocal models.Merchant
	var product models.Product
	var wg sync.WaitGroup
	var countProduct uint64

	fileXML, err := os.Open(pathFile)
	if err != nil {
		return
	}
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
			if start.Name.Local == "product" {
				dec.DecodeElement(&product, &start)
				countProduct++
			}
		}
	}

	log.Println("├─⇢ Product Counter:", countProduct)

	fileXML, err = os.Open(pathFile)
	if err != nil {
		return
	}
	dec = xml.NewDecoder(fileXML)

	chDBP := make(chan databases.DB)
	chDBM := make(chan databases.DB)
	chIDs := make(chan []string)
	chErr := make(chan error)

	go func() {
		chDBM <- databases.OpenDB(clientDB, databases.DBNameMerchants)
	}()

	go func() {
		chDBP <- databases.OpenDB(clientDB, databases.DBNameProducts)
	}()

	go func() {
		productsRemoteID, err := products.GetIDs()
		chIDs <- productsRemoteID
		chErr <- err
	}()

	dbMerchants := <-chDBM
	dbProducts := <-chDBP
	productsRemoteID := <-chIDs

	err = <-chErr
	if err != nil {
		return
	}

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
			case "header":
				if err = dec.DecodeElement(&merchant, &start); err != nil {
					break
				}

				err = uploadMerchant(merchant, dbMerchants)
				if err != nil {
					break
				}

			case "product":
				//Save products

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
				go uploadProduct(product, dbProducts, merchantLocal, &totalProductsUpload, &totalProductsUpdated, productsRemoteID, &wg)
			}
		}
	}

	wg.Wait()

	os.Remove("data/rakuten/decompress/" + filename)
	os.Remove("data/rakuten/" + filename + ".gz")

	return
}

func uploadMerchant(merchant utils.Merchant, dbMerchants databases.DB) (err error) {
	var merchantLocal models.Merchant

	merchantLocal = models.Merchant{
		ID:   merchant.MerchantID,
		Name: merchant.MerchantName,
	}

	var result interface{}
	databases.ReadElement(dbMerchants, merchantLocal.ID, &result, models.OptionsDB{})

	if result != nil {
		var merchantRemote models.Merchant
		if err = mapstructure.Decode(result, &merchantRemote); err != nil {
			return
		}
		merchantRemote.ID = result.(map[string]interface{})["_id"].(string)

		if merchantRemote != merchantLocal {
			rev := result.(map[string]interface{})["_rev"].(string)

			_, err = databases.UpdateElement(dbMerchants, merchantLocal.ID, rev, merchantLocal)
			if err != nil {
				if Verbose {
					log.Printf("├──⇢+ Error: %s Merchant #%s\n", err.Error(), merchantLocal.ID)
				}
				return
			}
			if Verbose {
				log.Printf("├──⇢+ Success: Merchant #%s, Updated\n", merchantLocal.ID)
			}
		}
	} else {
		_, _, err = databases.CreateElement(dbMerchants, merchantLocal)
		if err != nil {
			if Verbose {
				log.Printf("├──⇢+ Error: %s Merchant #%s\n", err.Error(), merchantLocal.ID)
			}
			return
		}
		if Verbose {
			log.Printf("├──⇢+ Success: Merchant #%s, Uploaded\n", merchantLocal.ID)
		}
	}

	return
}

func uploadProduct(product models.Product, dbProducts databases.DB, merchantLocal models.Merchant, totalProductsUpload, totalProductsUpdated *uint64, productsRemoteID []string, wg *sync.WaitGroup) (err error) {
	defer wg.Done()

	for {
		product.Merchant = merchantLocal

		var exists bool
		var productRemote models.Product
		var productRemoteREV string
		for _, v := range productsRemoteID {
			if v == product.ID {

				var result interface{}
				err = databases.ReadElement(dbProducts, product.ID, &result, models.OptionsDB{})
				if err != nil {
					return

				}
				if result != nil {
					exists = true

					productRemoteREV = result.(map[string]interface{})["_rev"].(string)
					if err = mapstructure.Decode(result, &productRemote); err != nil {
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
			continue
		}

		if exists {
			if productRemote != product {
				_, err = databases.UpdateElement(dbProducts, product.ID, productRemoteREV, product)
				if err != nil {
					if Verbose {
						log.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), product.ID)
					}
					continue
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
				continue
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
