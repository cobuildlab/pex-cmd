package main

import (
	"fmt"
	"log"
	"os"

	merchants "github.com/cobuildlab/pex-cmd/merchant-files"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pex-cmd",
	Short: "pex-retail-api Command line tool",
	Long:  "pex-retail-api Command line tool to perform automated tasks",
	Args:  cobra.MinimumNArgs(1),
}

func main() {
	log.SetPrefix(
		fmt.Sprintf("[%d] - ", os.Getpid()),
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rootCmd.AddCommand(merchants.RootMFCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
