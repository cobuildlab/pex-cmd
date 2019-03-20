package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountFiles(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		expectCount       uint64
		expectErrMessages string
	}{
		{
			"CountFilesExists",
			filepath.Join("..", "data"),
			1,
			"",
		},
		{
			"CountFilesNotExists",
			filepath.Join("..", "data", "rakuten-not-exists"),
			0,
			"An error has occurred counting local files",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			count, err := CountFiles(test.path, "")
			if err != nil {
				assert.Equal(t, test.expectErrMessages, err.Error())
			}

			assert.Equal(t, test.expectCount, count)

		})
	}
}
