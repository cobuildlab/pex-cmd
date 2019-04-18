package utils

import (
	"testing"

	"github.com/cobuildlab/pex-cmd/utils"
	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
)

func TestGetConnectionFTP(t *testing.T) {
	tests := []struct {
		name, host, port, username, password, expectErrMessages string
	}{
		{
			"ConnectionSuccessful",
			utils.FTPHost,
			utils.FTPPort,
			utils.FTPUsername,
			utils.FTPPassword,
			"",
		},
		{
			"UsernamePasswordIncorrect",
			utils.FTPHost,
			utils.FTPPort,
			"anonymous",
			"1010",
			"Incorrect username or password",
		},
		{
			"HostPortIncorrect",
			"anonymous",
			"1010",
			utils.FTPUsername,
			utils.FTPPassword,
			"Error making the connection",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, err := GetConnectionFTP(test.host, test.port, test.username, test.password)

			if err != nil {
				assert.Equal(t, test.expectErrMessages, err.Error())
			} else {
				assert.IsType(t, &ftp.ServerConn{}, client)
			}
		})
	}
}
