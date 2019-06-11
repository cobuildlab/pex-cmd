package merchants

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cobuildlab/pex-cmd/databases"
	"github.com/cobuildlab/pex-cmd/errors"
	"github.com/cobuildlab/pex-cmd/models"
	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/mitchellh/mapstructure"
)

const (
	//MerchantFileCompressExt Extension of compressed merchants files
	MerchantFileCompressExt = ".gz"

	//MerchantFileFormatExt Extension of the merchants files
	MerchantFileFormatExt = ".xml"

	//DecompressDir Name of decompression folder
	DecompressDir = "decompress"
)

var (
	//DecompressPath Decompression route
	DecompressPath = filepath.Join(utils.FTPPathFiles, "decompress")

	//FTPWildcardFilesFilter WildCard to filter merchants files from the FTP server
	FTPWildcardFilesFilter = fmt.Sprintf("*_%s_mp%s%s", utils.FTPSID, MerchantFileFormatExt, MerchantFileCompressExt)

	//FilesFilter Merchant files filter
	FilesFilter = fmt.Sprintf("_%s_mp%s%s", utils.FTPSID, MerchantFileFormatExt, MerchantFileCompressExt)
)

//CountProductsInMerchantFile Count the products within a merchant file
func CountProductsInMerchantFile(mf *os.File) (count int, err error) {
	var existNumberOfProducts bool

	dec := xml.NewDecoder(mf)

ForCountProduct:
	for {
		token, err := dec.Token()
		if err != nil {
			break
		} else if err == io.EOF {
			err = nil
			break
		}

		switch token.(type) {
		case xml.StartElement:
			start := token.(xml.StartElement)
			if start.Name.Local == "numberOfProducts" {
				dec.DecodeElement(&count, &start)

				existNumberOfProducts = true
				break ForCountProduct
			}
		}
	}

	if !existNumberOfProducts {
		err = errors.ErrorNumberOfProductsNotExist
		_, fileName := filepath.Split(mf.Name())
		exist, _ := utils.CheckExistence("data/rakuten/decompress/damaged")
		if !exist {
			os.Mkdir("data/rakuten/decompress/damaged", 0777)
		}
		os.Remove("data/rakuten/" + fileName + ".gz")
		os.Rename("data/rakuten/decompress/"+fileName, "data/rakuten/decompress/damaged/"+fileName)
	}

	return
}

//UploadMerchantInMerchantFile Upload the merchant of a Merchant file
func UploadMerchantInMerchantFile(mf *os.File, dbMerchants databases.DB) (merchant utils.Merchant, err error) {
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
