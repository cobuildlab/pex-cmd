package merchants

import (
	"fmt"
	"path/filepath"

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
