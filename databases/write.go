package databases

import (
	"time"

	"github.com/cobuildlab/pex-cmd/errors"
	"github.com/cobuildlab/pex-cmd/queue"
)

//QueueWriter ...
var QueueWriter = queue.NewClock(time.Second*1, uint(DBMaxWriting))

//CreateElement Create an item in the database
func CreateElement(db DB, element interface{}) (id string, rev string, err error) {
	QueueWriter.Add(1)
	defer QueueWriter.Done(1)

	id, rev, err = db.CreateDocument(element)
	//if err != nil {
	//	err = errors.ErrorCreateDocument
	//}

	return
}

//UpdateElement Update an item in the database
func UpdateElement(db DB, id string, rev string, element interface{}) (newRev string, err error) {
	QueueWriter.Add(1)
	defer QueueWriter.Done(1)

	newRev, err = db.UpdateDocument(id, rev, element)
	//if err != nil {
	//	err = errors.ErrorUpdateDocument
	//}

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
