package merchants

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/4geeks/pex-cmd/databases"
	"github.com/4geeks/pex-cmd/models"
	"github.com/4geeks/pex-cmd/services/products"
	"github.com/4geeks/pex-cmd/utils"
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
			fmt.Println("╭─Uploading", v, "...", fmt.Sprintf("%d/%d", i+1, len(args)))

			start := time.Now()

			totalProductsUpload, totalProductsUpdated, err := UploadFile(v, Verbose)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("├─⇢ ...", "Uploaded!", v)
			fmt.Println("│")
			fmt.Println("├──⇢ Uploaded products:", totalProductsUpload)
			fmt.Println("├──⇢ Updated products:", totalProductsUpdated)
			fmt.Println("╰──⇢ Duration:", time.Since(start))
			fmt.Println()
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

	fmt.Println("├─⇢ Product Counter:", countProduct)

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

				merchantLocal = models.Merchant{
					ID:   merchant.MerchantID,
					Name: merchant.MerchantName,
				}

				var result interface{}
				databases.ReadElement(dbMerchants, merchantLocal.ID, &result, models.OptionsDB{})

				if result != nil {
					var merchantRemote models.Merchant
					if err = mapstructure.Decode(result, &merchantRemote); err != nil {
						break
					}
					merchantRemote.ID = result.(map[string]interface{})["_id"].(string)

					if merchantRemote != merchantLocal {
						rev := result.(map[string]interface{})["_rev"].(string)

						_, err = databases.UpdateElement(dbMerchants, merchantLocal.ID, rev, merchantLocal)
						if err != nil {
							if Verbose {
								fmt.Printf("├──⇢+ Error: %s Merchant #%s\n", err.Error(), merchantLocal.ID)
							}
							break
						}
						if Verbose {
							fmt.Printf("├──⇢+ Success: Merchant #%s, Updated\n", merchantLocal.ID)
						}
					}
				} else {
					_, _, err = databases.CreateElement(dbMerchants, merchantLocal)
					if err != nil {
						if Verbose {
							fmt.Printf("├──⇢+ Error: %s Merchant #%s\n", err.Error(), merchantLocal.ID)
						}
						break
					}
					if Verbose {
						fmt.Printf("├──⇢+ Success: Merchant #%s, Uploaded\n", merchantLocal.ID)
					}
				}

			case "product":
				//Save products

				if err = dec.DecodeElement(&product, &start); err != nil {
					break
				}

				if product.Price.Currency != "USD" || product.Discount.Currency != "USD" {
					break
				}

				wg.Add(1)

				queueUpload <- true
				go func(wg *sync.WaitGroup, queueUpload <-chan bool, db databases.DB, productLocal models.Product, productsRemoteID []string) {
					defer func(wg *sync.WaitGroup) {
						wg.Done()
						<-queueUpload
					}(wg)

					for {
						productLocal.Merchant = merchantLocal

						var exists bool
						var productRemote models.Product
						var productRemoteREV string
						for _, v := range productsRemoteID {
							if v == productLocal.ID {

								var result interface{}
								err = databases.ReadElement(db, productLocal.ID, &result, models.OptionsDB{})
								if err != nil {
									break
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
								fmt.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), productLocal.ID)
							}
							continue
						}

						if exists {
							if productRemote != productLocal {
								_, err = databases.UpdateElement(db, productLocal.ID, productRemoteREV, productLocal)
								if err != nil {
									if Verbose {
										fmt.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), productLocal.ID)
									}
									continue
								}
								if Verbose {
									fmt.Printf("├──⇢+ Success: Product #%s, Updated\n", productLocal.ID)
								}
								atomic.AddUint64(&totalProductsUpdated, 1)
							}
						} else {
							_, _, err = databases.CreateElement(db, productLocal)
							if err != nil {
								if Verbose {
									fmt.Printf("├──⇢+ Error: %s Product #%s, Retrying\n", err.Error(), productLocal.ID)
								}
								continue
							}
							if Verbose {
								fmt.Printf("├──⇢+ Success: Product #%s, Uploaded\n", productLocal.ID)
							}
							atomic.AddUint64(&totalProductsUpload, 1)
						}
						break
					}

				}(&wg, queueUpload, dbProducts, product, productsRemoteID)

			}

			if err != nil {
				break
			}
		}

		if err != nil {
			break
		}
	}

	if err != nil {
		return
	}
	wg.Wait()

	return
}
