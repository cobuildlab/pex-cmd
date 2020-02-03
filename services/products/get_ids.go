package products

import (
	"fmt"
	"sort"

	//IbmCloudant "github.com/IBM-Cloud/go-cloudant"
	cloudant "github.com/cloudant-labs/go-cloudant"
	"github.com/cobuildlab/pex-cmd/databases"
	"github.com/cobuildlab/pex-cmd/models"
	"log"
	"time"
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
		//TODO: Check if the latest revision of the document can be query here to avoid doing the query on the uploadDocument coroutine
		Fields: []string{"_id"},
	}

	//TEST
	//resultData := make(map[string]string)
	//err = db.GetDocument("11732972158", &resultData, IbmCloudant.Options{})
	//log.Println("TEST-DEBUG:", resultData);
	//END TEST
	//os.Exit(1)

	startT := time.Now()
	log.Println("GetIDs:start time:", startT);
	result, err := databases.SearchElement(db, queryDB)
	if err != nil {
		return
	}
	log.Println("GetIDs:duration:", time.Since(startT));
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
	fmt.Println("GetIDs:productIDList:len:", len(productIDList))
	return
}

//GetIDs Get all product IDs
func GetProductsList() (RemoteProducts []models.RemoteProduct, err error) {
	client, err := cloudant.CreateClient(databases.Username, databases.Password, databases.Host, 5)
	if err != nil {
		log.Println("GetProductsList:duration:Error:", err.Error());
		return
	}
	db, err := client.GetOrCreate(databases.DBNameProducts)
	if err != nil {
		log.Println("GetProductsList:duration:Error:", err.Error());
		return
	}

	startT := time.Now()
	log.Println("GetProductsList:start time:", startT);
	q := cloudant.NewAllDocsQuery().
		Build()
	rows, err := db.All(q)

	if err != nil {
		log.Println("GetProductsList:duration:Error:", err.Error());
		return
	}
	log.Println("GetProductsList:duration:", time.Since(startT));
	if err != nil {
		log.Println("GetProductsList:duration:Error:", err.Error());
		return
	}
	for {
		row, more := <-rows
		if more {
			id := row.ID
			rev := row.Value.Rev
			product := models.RemoteProduct{id, id, rev}
			RemoteProducts = append(RemoteProducts, product)
		} else {
			break
		}
	}

	fmt.Println("GetProductsList:productIDList:len:", len(RemoteProducts))
	return
}

//GetIDs Get all product IDs
func GetSortedProductsList() (RemoteProducts []models.RemoteProduct, err error) {
	client, err := cloudant.CreateClient(databases.Username, databases.Password, databases.Host, 5)
	if err != nil {
		log.Println("GetSortedProductsList:duration:Error:", err.Error());
		return
	}
	db, err := client.GetOrCreate(databases.DBNameProducts)
	if err != nil {
		log.Println("GetSortedProductsList:duration:Error:", err.Error());
		return
	}

	startT := time.Now()
	log.Println("GetSortedProductsList:start time:", startT);
	q := cloudant.NewAllDocsQuery().
		//Limit(50).
		Build()
	rows, err := db.All(q)

	if err != nil {
		log.Println("GetSortedProductsList:Error:", err.Error());
		return
	}
	log.Println("GetSortedProductsList:duration:", time.Since(startT));
	if err != nil {
		log.Println("GetSortedProductsList:duration:Error:", err.Error());
		return
	}
	for {
		row, more := <-rows
		if more {
			id := row.ID
			rev := row.Value.Rev
			product := models.RemoteProduct{id, id, rev}
			RemoteProducts = append(RemoteProducts, product)
		} else {
			break
		}
	}

	sort.Slice(RemoteProducts, func(i, j int) bool {
		return RemoteProducts[i].ID > RemoteProducts[j].ID
	})

	fmt.Println("GetProductsList:productIDList:len:", len(RemoteProducts))
	return
}
