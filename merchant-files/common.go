package merchants

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cobuildlab/pex-cmd/utils"
)

const (
	//MerchantFileCompressExt Extension of compressed merchants files
	MerchantFileCompressExt = ".gz"

	//MerchantFileFormatExt Extension of the merchants files
	MerchantFileFormatExt = ".xml"

	//DecompressDir Name of decompression folder
	DecompressDir = "decompress"
)

var (
	//DecompressPath Decompression route
	DecompressPath = filepath.Join(utils.FTPPathFiles, "decompress")

	//FTPWildcardFilesFilter WildCard to filter merchants files from the FTP server
	FTPWildcardFilesFilter = fmt.Sprintf("*_%s_mp%s%s", utils.FTPSID, MerchantFileFormatExt, MerchantFileCompressExt)

	//FilesFilter Merchant files filter
	FilesFilter = fmt.Sprintf("_%s_mp%s%s", utils.FTPSID, MerchantFileFormatExt, MerchantFileCompressExt)
)

//TimeNow ...
func TimeNow() (now string) {
	dayN, minN, secN := time.Now().Clock()
	day, min, sec := strconv.Itoa(dayN), strconv.Itoa(minN), strconv.Itoa(secN)

	if len(day) < 2 {
		day = "0" + day
	}
	if len(min) < 2 {
		min = "0" + min
	}
	if len(sec) < 2 {
		sec = "0" + sec
	}

	now = fmt.Sprintf(
		"%s:%s:%s - ",
		day,
		min,
		sec,
	)

	return
}
