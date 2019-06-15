package merchants

import (
	"github.com/spf13/cobra"
)

//Verbose Verbose mode
var Verbose bool

//LimitSize File size limit
var LimitSize uint64

func init() {
	CmdDownloadAll.Flags().Uint64VarP(&LimitSize, "limitsize", "l", 9999999999, "Limit file size")

	CmdDownload.AddCommand(CmdDownloadAll)
	CmdDownload.AddCommand(CmdDownloadFile)
	CmdDownload.AddCommand(CmdDownloadList)

	CmdUploadFile.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")
	CmdUploadAll.Flags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose output")

	CmdUpload.AddCommand(CmdUploadList)
	CmdUpload.AddCommand(CmdUploadCount)
	CmdUpload.AddCommand(CmdUploadFile)
	CmdUpload.AddCommand(CmdUploadAll)

	RootMFCmd.AddCommand(CmdDownload)
	RootMFCmd.AddCommand(CmdUpload)
}

//RootMFCmd Root commands mf
var RootMFCmd = &cobra.Command{
	Use:   "mf",
	Short: "Operations on merchants files",
	Long:  "Operations on merchants files, Download and upload",
	Args:  cobra.MinimumNArgs(1),
}

//CmdDownload Subcommands download
var CmdDownload = &cobra.Command{
	Use:   "download",
	Short: "Merchants files download operations",
	Long:  "Merchants files download operations, FTP server download merchant files",
	Args:  cobra.MinimumNArgs(1),
}

//CmdUpload Subcommands upload
var CmdUpload = &cobra.Command{
	Use:   "upload",
	Short: "Merchants files upload operations",
	Long:  "Merchants files upload operations, Merchants files available to upload to the database",
	Args:  cobra.MinimumNArgs(1),
}
