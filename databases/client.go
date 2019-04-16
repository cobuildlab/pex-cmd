package databases

import (
	cloudant "github.com/IBM-Cloud/go-cloudant"
	"github.com/cobuildlab/pex-cmd/errors"
)

//NewClient Returns the interface of a database client
func NewClient(username string, password string) (client Client, err error) {
	client, err = cloudant.NewClient(username, password)
	if err != nil {
		err = errors.ErrorConnectingDatabase
	}

	return
}
