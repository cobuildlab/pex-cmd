package utils

import (
	"io/ioutil"
	"strings"

	"github.com/cobuildlab/pex-cmd/errors"
)

//CountFiles Count the files in a specified path
func CountFiles(path, contains string) (countLocalFiles uint64, err error) {
	localFilesDir, err := ioutil.ReadDir(path)
	if err != nil {
		err = errors.ErrorReadingData

		return
	}

	for _, file := range localFilesDir {
		if strings.Contains(file.Name(), contains) {
			countLocalFiles++
		}
	}

	return
}
