package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/4geeks/pex-cmd/errors"

	"github.com/jlaffaye/ftp"
)

//GetConnectionFTP Returns a connection to the FTP server
func GetConnectionFTP(host, port, username, password string) (client *ftp.ServerConn, err error) {
	client, err = ftp.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		err = errors.ErrorConnection

		return
	}

	if err = client.Login(username, password); err != nil {
		err = errors.ErrorIncorrectCredentials

		return
	}

	return
}

//DownloadGzipFileFTP Download a gzip file from an FTP entry in the specified path
func DownloadGzipFileFTP(entryName string, path string) (err error) {
	client, err := GetConnectionFTP(
		FTPHost, FTPPort,
		FTPUsername, FTPPassword,
	)
	if err != nil {
		return
	}
	defer client.Quit()
	defer client.Logout()

	filePath := filepath.Join(path, entryName)

	exist, err := CheckExistence(filePath)
	if err != nil {
		return
	}

	if exist {
		ok, _ := CheckIntegrityGzip(filePath)
		if ok {
			return
		}
	}

	for {
		var dstFile *os.File
		dstFile, err = os.Create(filePath)
		if err != nil {
			err = errors.ErrorCreatingFileDirectory

			return
		}

		var reader *ftp.Response
		reader, err = client.Retr(entryName)
		if err != nil {
			err = errors.ErrorFTPGetFile

			return
		}

		_, err = io.Copy(dstFile, reader)
		if err != nil {
			dstFile.Close()
			reader.Close()

			err = errors.ErrorFTPRecordingData

			return
		}
		dstFile.Close()
		reader.Close()

		if ok, _ := CheckIntegrityGzip(filePath); ok {
			break
		}
	}

	return
}
