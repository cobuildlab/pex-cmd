package models

import cloudant "github.com/IBM-Cloud/go-cloudant"

//OptionsDB Options for consultation in the database
type OptionsDB cloudant.Options

//SelectorDB Selector to query the database
type SelectorDB map[string]interface{}

//QueryDB Query for the database
type QueryDB cloudant.Query

//AllDocsResult Model to obtain all the elements of a database
type AllDocsResult struct {
	TotalRows int                      `json:"total_rows"`
	Offset    int                      `json:"offset"`
	Rows      []map[string]interface{} `json:"rows"`
}
