package utils

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cobuildlab/pex-cmd/errors"
)

//CheckIntegrityGzip Verify the integrity of a compressed file with gzip
func CheckIntegrityGzip(path string) (ok bool, err error) {
	_, fileInputName := filepath.Split(path)

	fileExt := filepath.Ext(fileInputName)

	if fileExt != ".gz" {
		err = errors.ErrorInformationNotProvided

		return
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.ErrorReadingData

		return
	}

	defer file.Close()

	compress, err := gzip.NewReader(file)
	if err != nil {
		err = errors.ErrorDecompressingFile

		return
	}
	defer compress.Close()

	fileTemp, err := ioutil.TempFile(FTPPathFiles, "decompressCheck")
	if err != nil {
		return
	}
	defer os.Remove(fileTemp.Name())

	_, err = io.Copy(fileTemp, compress)
	if err != nil {
		err = errors.ErrorReadingData
		return
	}

	ok = true

	return
}
