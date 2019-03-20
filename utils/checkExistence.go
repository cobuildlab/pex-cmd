package utils

import (
	"os"

	"github.com/4geeks/pex-cmd/errors"
)

//CheckExistence Verify the existence of a file or directory
func CheckExistence(path string) (exist bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		err = errors.ErrorCheckExistence
		return
	}

	exist = true

	return
}
