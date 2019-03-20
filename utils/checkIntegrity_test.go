package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckIntegrityGzip(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		expectOk          bool
		expectErrMessages string
	}{
		{
			"GzipIntegrated",
			filepath.Join("..", "data", "rakuten", "1110_3551736_mp.xml.gz"),
			true,
			"",
		},
		{
			"FileReadingError",
			filepath.Join("..", "data", "rakuten", "decompress", "13923_3551736_mp.xml.gz"),
			false,
			"An error has occurred reading the local file",
		},
		{
			"FileNotGzip",
			filepath.Join("..", "data", "rakuten", "decompress", "1110_3551736_mp.xml"),
			false,
			"A decompression error has occurred, A GZIP file (.gz) was expected",
		},
		{
			"GzipNotIntegrated",
			filepath.Join("..", "data", "rakuten", "147_3551736_mp.xml.gz"),
			false,
			"An error has occurred decompressing file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := CheckIntegrityGzip(test.path)

			assert.Equal(t, test.expectOk, ok)

			if !test.expectOk {
				assert.Equal(t, test.expectErrMessages, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
