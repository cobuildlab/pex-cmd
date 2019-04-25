package main

import (
	"fmt"
	"log"
	"os"

	merchants "github.com/cobuildlab/pex-cmd/merchant-files"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pex-cmd",
	Short: "pex-retail-api Command line tool",
	Long:  "pex-retail-api Command line tool to perform automated tasks",
	Args:  cobra.MinimumNArgs(1),
}

func main() {
	f, err := os.OpenFile("pex.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.SetOutput(f)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rootCmd.AddCommand(merchants.RootMFCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
