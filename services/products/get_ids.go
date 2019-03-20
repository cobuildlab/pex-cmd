package products

import (
	"github.com/4geeks/pex-cmd/databases"
	"github.com/4geeks/pex-cmd/models"
)

//GetIDs Get all product IDs
func GetIDs() (productIDList []string, err error) {
	clientDB, err := databases.NewClient(
		databases.Username,
		databases.Password,
	)
	if err != nil {
		return
	}
	db := databases.OpenDB(clientDB, databases.DBNameProducts)

	queryDB := models.QueryDB{
		Selector: models.SelectorDB{},
		Fields:   []string{"_id"},
	}

	result, err := databases.SearchElement(db, queryDB)
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	for _, element := range result {
		id := element.(map[string]interface{})["_id"]
		if id == nil {
			continue
		}
		productIDList = append(productIDList, id.(string))
	}

	return
}
