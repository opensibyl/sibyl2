package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "sibyl server cmd",
	Long:  `sibyl server cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Execute(server.DefaultExecuteConfig())
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
