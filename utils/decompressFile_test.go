package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompressFileGzip(t *testing.T) {
	tests := []struct {
		name              string
		fileInput         string
		pathOutput        string
		expectErrMessages string
	}{
		{
			"DecompressFileSuccessful",
			filepath.Join("..", "data", "rakuten", "1110_3551736_mp.xml.gz"),
			filepath.Join("..", "data", "rakuten", "decompress"),
			"",
		},
		{
			"FileNotGzip",
			filepath.Join("..", "data", "rakuten", "decompress", "1110_3551736_mp.xml"),
			filepath.Join("..", "data", "rakuten", "decompress"),
			"A decompression error has occurred, A GZIP file (.gz) was expected",
		},
		{
			"GzipNotIntegrated",
			filepath.Join("..", "data", "rakuten", "13923_3551736_mp.xml.gz"),
			filepath.Join("..", "data", "rakuten", "decompress"),
			"GZIP file provided is not integrated",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := DecompressFileGzip(test.fileInput, test.pathOutput)
			if err != nil {
				assert.Equal(t, test.expectErrMessages, err.Error())
			}
		})
	}
}
