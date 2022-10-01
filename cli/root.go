package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "sibyl",
	Short: "sibyl cli",
	Long:  `sibyl cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root cmd from sibyl 2")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
