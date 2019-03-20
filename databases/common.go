package databases

import (
	"os"
	"strconv"

	//Autoload .env
	_ "github.com/joho/godotenv/autoload"
)

var (
	//Username Username of the database
	Username = os.Getenv("CLOUDANT_USERNAME")

	//Password Password of the database
	Password = os.Getenv("CLOUDANT_PASSWORD")

	//Host Host of database
	Host = "https://" + Username + ".cloudant.com"

	//DBNameTest Database test
	DBNameTest = os.Getenv("CLOUDANT_DATABASE_TEST")

	//DBNameMerchants Database merchants
	DBNameMerchants = os.Getenv("CLOUDANT_DATABASE_MERCHANTS")

	//DBNameProducts Database products
	DBNameProducts = os.Getenv("CLOUDANT_DATABASE_PRODUCTS")

	//DBNameCategories Database categories
	DBNameCategories = os.Getenv("CLOUDANT_DATABASE_CATEGORIES")

	//DBNameUsers Database users
	DBNameUsers = os.Getenv("CLOUDANT_DATABASE_USERS")

	//DBNameStatistics Database statistics
	DBNameStatistics = os.Getenv("CLOUDANT_DATABASE_STATISTICS")

	//DBNamePartners Database Partners
	DBNamePartners = os.Getenv("CLOUDANT_DATABASE_PARTNERS")

	//DBMaxWriting Maximum possible writes in the database per second
	DBMaxWriting, _ = strconv.Atoi(os.Getenv("CLOUDANT_DATABASE_MAX_WRITING"))

	//DBMaxReading Maximum possible readings in the database per second
	DBMaxReading, _ = strconv.Atoi(os.Getenv("CLOUDANT_DATABASE_MAX_READING"))
)
