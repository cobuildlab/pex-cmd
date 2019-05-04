package databases

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	cloudant "github.com/IBM-Cloud/go-cloudant"
	queue "github.com/arthurnavah/go-queue"
	"github.com/cobuildlab/pex-cmd/errors"
	"github.com/cobuildlab/pex-cmd/models"
	couchdb "github.com/timjacobi/go-couchdb"
)

var queueRead = queue.NewClock(time.Second*1, uint(DBMaxReading))

//ReadAllElements Read all the elements of a database
func ReadAllElements(db DB, result interface{}, opts models.OptionsDB) (err error) {
	queueRead.Add(1)
	defer queueRead.Done(1)

	err = db.GetAllDocument(result, cloudant.Options(opts))
	if err != nil {
		err = errors.ErrorGetAllDocument
	}

	return
}

//ReadElement Read element of a database
func ReadElement(db DB, id string, doc interface{}, opts models.OptionsDB) (err error) {
	queueRead.Add(1)
	defer queueRead.Done(1)

	err = db.GetDocument(id, doc, cloudant.Options(opts))
	if err != nil {
		switch err.(type) {
		case *couchdb.Error:
			errC := err.(*couchdb.Error)
			if errC.StatusCode == 404 {
				err, doc = nil, nil
				return
			}
		}
		err = errors.ErrorGetDocument
	}

	return
}

//SearchElement Search element of a database
func SearchElement(db DB, query models.QueryDB) (result []interface{}, err error) {
	queueRead.Add(1)
	defer queueRead.Done(1)

	result, err = db.SearchDocument(cloudant.Query(query))
	if err != nil {
		err = errors.ErrorSearchDocument
	}

	return
}

//SearchDesignDocument Perform a search using a Design Document, Return a SearchResp with the response of the database
func SearchDesignDocument(db, name, index, query string, page, limit, maxLimit int) (result cloudant.SearchResp, err error) {
	queueRead.Add(1)
	defer queueRead.Done(1)

	if page == 0 {
		page = 1
	}

	path := "/_design/" + name + "/_search/" + index

	send := map[string]interface{}{
		"query": query,
		"limit": maxLimit,
	}

	endPoint := Host + "/" + db + path

	sendJ, _ := json.Marshal(send)
	req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(sendJ))
	if err != nil {
		return
	}
	req.SetBasicAuth(Username, Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyN, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var response cloudant.SearchResp
	json.Unmarshal(bodyN, &response)

	result.Bookmark = response.Bookmark
	result.Num = response.Num

	if page >= 0 {
		if len(response.Rows) >= page*limit+limit {
			result.Rows = response.Rows[(page*limit)-1 : (page*limit+limit)-1]
		} else if len(response.Rows) >= limit-1 {
			result.Rows = response.Rows[:limit-1]
		} else {
			result.Rows = response.Rows
		}
	} else {
		result.Rows = response.Rows
	}

	return
}
