package databases

import (
	cloudant "github.com/IBM-Cloud/go-cloudant"
)

//DB Interface of the database
type DB interface {
	GetAllDocument(result interface{}, opts cloudant.Options) error

	CreateDocument(doc interface{}) (string, string, error)

	GetDocument(id string, doc interface{}, opts cloudant.Options) error

	SearchDocument(query cloudant.Query) (result []interface{}, err error)

	UpdateDocument(id, rev string, doc interface{}) (string, error)
}

//Client Interface database client
type Client interface {
	CreateDB(dbName string) (*cloudant.DB, error)

	DB(name string) *cloudant.DB

	DeleteDB(dbName string) error
}
