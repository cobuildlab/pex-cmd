package models

//QueryProduct Search for products search
type QueryProduct struct {
	Page  int `json:"page,omitempty" form:"page" query:"page"`
	Limit int `json:"limit,omitempty" form:"limit" query:"limit"`

	Keywords string  `json:"keywords,omitempty" form:"keywords" query:"keywords"`
	Brand    string  `json:"brand,omitempty" form:"brand" query:"brand"`
	Category string  `json:"category,omitempty" form:"category" query:"category"`
	MinPrice float64 `json:"minPrice,omitempty" form:"minPrice" query:"minPrice"`
	MaxPrice float64 `json:"maxPrice,omitempty" form:"maxPrice" query:"maxPrice"`
}

//Product Product model
type Product struct {
	ID   string `xml:"product_id,attr" json:"_id,omitempty"`
	Name string `xml:"name,attr" json:"name,omitempty"`

	SkuNumber    string `xml:"sku_number,attr" json:"skuNumber,omitempty"`
	Manufacturer string `xml:"manufacturer_name,attr" json:"manufacturer,omitempty"`
	PartNumber   string `xml:"part_number,attr" json:"partNumber,omitempty"`

	Merchant Merchant `json:"merchant,omitempty" xml:"merchant"`

	Category struct {
		Primary   string `xml:"primary" json:"primary,omitempty"`
		Secondary string `xml:"secondary" json:"secondary,omitempty"`
	} `xml:"category" json:"category,omitempty"`

	URL struct {
		Product string `xml:"product" json:"product,omitempty"`
		Image   string `xml:"productImage" json:"image,omitempty"`
		Buy     string `xml:"buyLink" json:"buy,omitempty"`
	} `xml:"URL" json:"url,omitempty"`

	Description struct {
		Short string `xml:"short" json:"short,omitempty"`
		Long  string `xml:"long" json:"long,omitempty"`
	} `xml:"description" json:"description,omitempty"`

	Discount struct {
		Currency string  `xml:"currency,attr" json:"currency,omitempty"`
		Type     string  `xml:"type" json:"type,omitempty"`
		Amount   float64 `xml:"amount" json:"amount,omitempty"`
	} `xml:"discount" json:"discount,omitempty"`

	Price struct {
		Currency  string  `xml:"currency,attr" json:"currency,omitempty"`
		Sale      float64 `xml:"sale" json:"sale,omitempty"`
		BeginData string  `xml:"begin_data" json:"beginData,omitempty"`
		EndData   string  `xml:"end_data" json:"endData,omitempty"`
		Retail    float64 `xml:"retail" json:"retail,omitempty"`
	} `xml:"price" json:"price,omitempty"`

	Brand string `xml:"brand" json:"brand,omitempty"`

	Shipping struct {
		Availability string `xml:"availability" json:"availability,omitempty"`
		Information  string `xml:"information" json:"information,omitempty"`

		Cost struct {
			Currency string  `xml:"currency" json:"currency,omitempty"`
			Amount   float64 `xml:"amount" json:"amount,omitempty"`
		} `xml:"cost" json:"cost,omitempty"`
	} `xml:"shipping" json:"shipping,omitempty"`

	Keywords string `xml:"keywords" json:"keywords,omitempty"`

	UPC string `xml:"upc" json:"upc,omitempty"`
	M1  string `xml:"m1" json:"m1,omitempty"`

	Pixel string `xml:"pixel" json:"pixel,omitempty"`

	// // Attributes struct {}

	Modification string `xml:"modification" json:"modification,omitempty"`
}
