package merchants

import (
	"encoding/xml"
	"fmt"
	cloudant "github.com/cloudant-labs/go-cloudant"
	"github.com/cobuildlab/pex-cmd/databases"
	"github.com/cobuildlab/pex-cmd/models"
	"github.com/cobuildlab/pex-cmd/services/products"
	"sort"

	//"github.com/cobuildlab/pex-cmd/services/products"
	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/mitchellh/mapstructure"
	"io"
	"os"
	"path/filepath"
	"time"

	//"fmt"
	"log"
	//"time"

	//"github.com/cobuildlab/pex-cmd/utils"
	"github.com/spf13/cobra"
)

//CmdUploadAll Command to upload all the merchant files to the database<Paste>
var UploadProductsCMD = &cobra.Command{
	Use:   "products",
	Short: "Upload all the merchants files to the database",
	Long:  "Upload all the merchants files to the database",
	Run: func(cmd *cobra.Command, args []string) {
		//f, err := os.OpenFile("upload.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		//if err != nil {
		//	log.Println(err)
		//	os.Exit(1)
		//}
		//log.SetOutput(f)

		fileList, err := UploadList()
		if err != nil {
			log.Println(err)
			return
		}

		for _, file := range (fileList) {
			log.Println(file.Name, file.Size)
		}

		//some, err := products.GetIDs()
		//log.Println("TEST:LIST:", some);
		remoteProductsList, err := products.GetSortedProductsList()
		if err != nil {
			return
		}

		//for _, remoteProduct := range (remoteProductsList) {
		//	log.Println(remoteProduct.ID, remoteProduct.Rev)
		//}

		//remoteProductsList := []models.RemoteProduct{}

		startT := time.Now()
		var totalProductsUploadT, totalProductsUpdatedT uint64
		for i := 0; len(fileList) != 0; i++ {
			//if i == 2 {
			//	os.Exit(1)
			//}

			v := fileList[0]
			start := time.Now()
			log.Println("├─⇢ ...", "Start Processing!", v.Name)
			totalProductsUpload, totalProductsUpdated, err := UploadProductFile(v.Name, remoteProductsList, Verbose)
			if err != nil {
				log.Println("├─⇢ ...", "ERROR!", err.Error())
			}

			log.Println("├─⇢ ...", "Uploaded!", v.Name)
			log.Println("│")
			log.Println("├──⇢ Uploaded products:", totalProductsUpload)
			log.Println("├──⇢ Updated products:", totalProductsUpdated)
			log.Println("╰──⇢ Duration:", time.Since(start))
			log.Println()

			totalProductsUploadT += totalProductsUpload
			totalProductsUpdatedT += totalProductsUpdated

			log.Println("├─⇢ ...", "Removing:", v.Name)
			os.Remove("data/rakuten/decompress/" + v.Name)
			os.Remove("data/rakuten/" + v.Name + ".gz")

			fileList, err = UploadList()
			//for _, file := range (fileList) {
			//	log.Println("ITERATION:", i, file.Name, file.Size)
			//}
			if err != nil {
				log.Println(err)
				return
			}
			//os.Exit(1)
		}

		fmt.Println()
		log.Println("╭──⇢ Total Uploaded products:", totalProductsUploadT)
		log.Println("├──⇢ Total Updated products:", totalProductsUpdatedT)
		log.Println("╰──⇢ Duration:", time.Since(startT))
		fmt.Println()
	},
}

//UploadFile Upload a merchant file to the database
func UploadProductFile(filename string, remotesProducts []models.RemoteProduct, verbose bool) (totalProductsUpload, totalProductsUpdated uint64, err error) {
	log.Println("UploadProductFile:", filename)
	var merchant utils.Merchant
	var product models.Product

	client, err := cloudant.CreateClient(databases.Username, databases.Password, databases.Host, 10)
	log.Println("UploadProductFile:before:GetOrCreate:", filename)
	db, err := client.GetOrCreate(databases.DBNameProducts)
	log.Println("UploadProductFile:before:Bulk:", filename)
	uploader := db.Bulk(500, 2097152, 120) // 2mb

	log.Println("├─⇢ Opening:", filename)
	pathFile := filepath.Join(DecompressPath, filename)
	fileXML, err := os.Open(pathFile)
	if err != nil {
		log.Println("├─⇢ Error Opening File:")
		return
	}

	// Creating the Merchant File if doesn't exist
	merchant, err = createOrUpdateMerchantFromProductsFile(fileXML)
	if err != nil {
		log.Println("├─⇢ Error OpenicreateOrUpdateMerchantFromProductsFile:", err.Error())
		return
	}
	fileXML.Close()

	fileXML, err = os.Open(pathFile)
	if err != nil {
		log.Println("├─⇢ Error Opening File:")
		return
	}
	defer fileXML.Close()

	// Loading Products
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
				if err = dec.DecodeElement(&product, &start); err != nil {
					break
				}

				if product.Price.Currency != "" && product.Price.Currency != "USD" {
					break
				}

				if product.Discount.Currency != "" && product.Discount.Currency != "USD" {
					break
				}
				findRevIfExist(&product, merchant, remotesProducts)
				uploader.Upload(product)
				log.Println("TEST DEBUG:", product.Rev)
				log.Println("TEST DEBUG:", product.Rev != "")

				if product.Rev != "" {
					totalProductsUpdated++
				} else {
					totalProductsUpload++
				}
			}
		}

	}

	uploader.Flush()
	return
}

// Upload a product to the database
func findRevIfExist(product *models.Product, merchantLocal utils.Merchant, remoteProducts []models.RemoteProduct) (err error) {
	log.Println("findRevIfExist:find:", product.ID)
	product.Merchant = models.Merchant{
		ID:   merchantLocal.MerchantID,
		Name: merchantLocal.MerchantName,
	}
	product.Rev = ""

	// Binary Search for sorted List
	remoteProductsLen := len(remoteProducts)
	i := sort.Search(remoteProductsLen, func(i int) bool { return remoteProducts[i].ID <= product.ID })
	if i < remoteProductsLen && remoteProducts[i].ID == product.ID {
		log.Println("findRevIfExist:MATCH FOUND:", product.ID)
	} else {
		log.Println("findRevIfExist:MATCH NOT FOUND: need to create:", product.ID)
	}

	return
}

func createOrUpdateMerchantFromProductsFile(mf *os.File) (merchant utils.Merchant, err error) {
	log.Println("createOrUpdateMerchantFromProductsFile:")

	clientDB, err := databases.NewClient(databases.Username, databases.Password)
	if err != nil {
		return
	}

	dbMerchants := databases.OpenDB(clientDB, databases.DBNameMerchants)

	dec := xml.NewDecoder(mf)

ForDecodeMerchant:
	for {
		token, err := dec.Token()
		if err == io.EOF {
			err = nil
			break
		}

		switch token.(type) {
		case xml.StartElement:
			start := token.(xml.StartElement)
			if start.Name.Local == "header" {
				dec.DecodeElement(&merchant, &start)
				break ForDecodeMerchant
			}
		}

	}

	merchantLocal := models.Merchant{
		ID:   merchant.MerchantID,
		Name: merchant.MerchantName,
	}
	log.Println("├─⇢ Merchant ID:", merchantLocal.ID)

	var result interface{}
	databases.ReadElement(dbMerchants, merchantLocal.ID, &result, models.OptionsDB{})

	if result != nil {
		var merchantRemote models.Merchant
		if err = mapstructure.Decode(result, &merchantRemote); err != nil {
			return
		}

		if _, ok := result.(map[string]interface{})["_id"].(string); ok {
			merchantRemote.ID = result.(map[string]interface{})["_id"].(string)
		}

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
