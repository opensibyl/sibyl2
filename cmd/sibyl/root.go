package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "sibyl",
	Short: "sibyl cmd",
	Long:  `sibyl cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root cmd from sibyl 2")
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
