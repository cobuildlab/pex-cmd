package databases

import (
	"time"

	"github.com/cobuildlab/pex-cmd/errors"
)

var queueWriter = make(chan bool, DBMaxWriting)
var queueWriterDone = make(chan bool, DBMaxWriting)

func init() {
	go func() {
		c := time.Tick(time.Second * 1)
		for range c {
			for i := 0; i < len(queueWriterDone); i++ {
				<-queueWriter
				<-queueWriterDone
			}
		}
	}()
}

//CreateElement Create an item in the database
func CreateElement(db DB, element interface{}) (id string, rev string, err error) {
	queueWriter <- true
	defer func() {
		queueWriterDone <- true
	}()

	id, rev, err = db.CreateDocument(element)
	if err != nil {
		err = errors.ErrorCreateDocument
	}

	return
}

//UpdateElement Update an item in the database
func UpdateElement(db DB, id string, rev string, element interface{}) (newRev string, err error) {
	queueWriter <- true
	defer func() {
		queueWriterDone <- true
	}()

	newRev, err = db.UpdateDocument(id, rev, element)
	if err != nil {
		err = errors.ErrorUpdateDocument
	}

	return
}

//DeleteDB Delete a database
func DeleteDB(client Client, name string) (err error) {
	err = client.DeleteDB(name)
	if err != nil {
		err = errors.ErrorDeleteDatabase
	}

	return
}
