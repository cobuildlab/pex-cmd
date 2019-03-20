package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestCheckExistence Test CheckExistence
func TestCheckExistence(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectExists bool
		expectErr    bool
	}{
		{
			"FileExists",
			filepath.Join("checkExistence.go"),
			true,
			false,
		},
		{
			"FileNotExists",
			filepath.Join("fileNotExists01.bin"),
			false,
			false,
		},
		{
			"DirExists",
			filepath.Join("..", "utils"),
			true,
			false,
		},
		{
			"DirNotExists",
			filepath.Join("..", "dirNotExists"),
			false,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists, err := CheckExistence(test.path)

			assert.Equal(t, test.expectExists, exists)

			if test.expectErr {
				assert.Error(t, err)
			}
		})

	}
}
