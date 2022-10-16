package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/extras/casedoctor"
	"log"
	"os"
)

var userSrc string
var userOutputFile string

var rootCmd = &cobra.Command{
	Use:   "casedoctor",
	Short: "casedoctor cmd",
	Long:  `casedoctor cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := casedoctor.CheckCases(userSrc)
		if err != nil {
			panic(err)
		}
		output, err := json.MarshalIndent(&result, "", "  ")
		if err != nil {
			panic(err)
		}

		if userOutputFile == "" {
			userOutputFile = fmt.Sprintf("casedoctor.json")
		}
		err = os.WriteFile(userOutputFile, output, 0644)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&userSrc, "src", ".", "src dir path")
	err := rootCmd.MarkPersistentFlagRequired("src")
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().StringVar(&userOutputFile, "output", "", "output file")
	if err != nil {
		panic(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	Execute()
}
