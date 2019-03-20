package utils

import (
	"time"

	"github.com/4geeks/pex-cmd/models"
)

//FileInfo Custom FileInfo
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

//MerchantFile Merchant file structure
type MerchantFile struct {
	Header Merchant `xml:"header" json:"header"`

	Products []models.Product `xml:"product" json:"products"`

	Trailer struct {
		NumberOfProducts int `xml:"numberOfProducts" json:"numberOfProducts"`
	} `xml:"trailer" json:"trailer"`
}

//Merchant Merchant model inside a merchant file
type Merchant struct {
	MerchantID   string `xml:"merchantId" json:"merchantID"`
	MerchantName string `xml:"merchantName" json:"merchantName"`
	CreatedOn    string `xml:"createdOn" json:"createdOn"`
}
