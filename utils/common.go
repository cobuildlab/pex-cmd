package utils

import (
	"os"

	//Autoload .env
	_ "github.com/joho/godotenv/autoload"
)

var (
	//FTPUsername FTP server's username
	FTPUsername = os.Getenv("FTP_USERNAME")

	//FTPPassword FTP server password
	FTPPassword = os.Getenv("FTP_PASSWORD")

	//FTPHost FTP server host
	FTPHost = os.Getenv("FTP_HOST")

	//FTPPort FTP server port
	FTPPort = os.Getenv("FTP_PORT")

	//FTPSID Rakuten SID for the FTP server
	FTPSID = os.Getenv("FTP_SID")

	//FTPPathFiles File path for the merchants files of the FTP server
	FTPPathFiles = os.Getenv("FTP_PATH_FILES")
)
