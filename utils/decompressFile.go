package utils

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cobuildlab/pex-cmd/errors"
)

//DecompressFileGzip Uncompress a gzip file on a specified path
func DecompressFileGzip(fileInput, pathOutput string) (err error) {
	_, fileInputName := filepath.Split(fileInput)

	fileExt := filepath.Ext(fileInputName)

	if fileExt != ".gz" {
		err = errors.ErrorInformationNotExpected

		return
	}

	ok, _ := CheckIntegrityGzip(fileInput)
	if !ok {
		err = errors.ErrorFileCorrupted

		return
	}

	fileName := strings.TrimSuffix(fileInputName, fileExt)

	fileOutputPath := filepath.Join(pathOutput, fileName)

	exist, _ := CheckExistence(pathOutput)
	if !exist {
		err = os.MkdirAll(pathOutput, 0777)
		if err != nil {
			err = errors.ErrorCreatingFileDirectory

			return
		}
	}

	exist, _ = CheckExistence(fileOutputPath)

	if !exist {

		var file *os.File
		file, err = os.Open(fileInput)
		if err != nil {
			err = errors.ErrorReadingData

			return
		}

		var decompress *gzip.Reader
		decompress, err = gzip.NewReader(file)
		if err != nil {
			err = errors.ErrorDecompressingFile
			return
		}

		defer decompress.Close()

		var dstFile *os.File
		dstFile, err = os.Create(fileOutputPath)
		if err != nil {
			err = errors.ErrorCreatingFileDirectory

			return
		}

		defer dstFile.Close()

		io.Copy(dstFile, decompress)
	}

	return
}
