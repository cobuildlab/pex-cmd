package databases

//OpenDB Returns a database
func OpenDB(client Client, name string) (db DB) {
	db, err := client.CreateDB(name)
	if err != nil {
		db = client.DB(name)
	}

	return db
}
