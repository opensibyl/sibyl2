package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/cmd/sibyl/extract"
	"github.com/williamfzc/sibyl2/cmd/sibyl/server"
	"github.com/williamfzc/sibyl2/cmd/sibyl/upload"
)

var rootCmd = &cobra.Command{
	Use:   "sibyl",
	Short: "sibyl cmd",
	Long:  `sibyl cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	Execute()
}

func init() {
	extractCmd := extract.NewExtractCmd()
	rootCmd.AddCommand(extractCmd)

	serverCmd := server.NewServerCmd()
	rootCmd.AddCommand(serverCmd)

	uploadCmd := upload.NewUploadCmd()
	rootCmd.AddCommand(uploadCmd)
}
