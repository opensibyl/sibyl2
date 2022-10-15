package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/extras/storytrack"
	"github.com/williamfzc/sibyl2/pkg/core"
	"log"
	"os"
)

var userSrc string
var userLangType string
var userIds []int
var userOutputFile string

var rootCmd = &cobra.Command{
	Use:   "storytrack",
	Short: "storytrack cmd",
	Long:  `storytrack cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		langType := core.LangTypeValueOf(userLangType)
		if langType == core.LangUnknown {
			panic(errors.New("unknown lang type: " + userLangType))
		}
		result, err := storytrack.TrackWithSharpId(userSrc, "", userIds, core.LangGo)
		if err != nil {
			panic(err)
		}
		output, err := json.MarshalIndent(&result, "", "  ")
		if err != nil {
			panic(err)
		}

		if userOutputFile == "" {
			userOutputFile = fmt.Sprintf("story-%v.json", userIds)
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

	rootCmd.PersistentFlags().StringVar(&userLangType, "lang", "", "lang type of your source code")
	err = rootCmd.MarkPersistentFlagRequired("lang")
	if err != nil {
		panic(err)
	}

	rootCmd.PersistentFlags().IntSliceVar(&userIds, "ids", []int{}, "story ids")
	err = rootCmd.MarkPersistentFlagRequired("ids")
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
