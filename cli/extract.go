package cli

import (
	"errors"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	"sibyl2/pkg/core"
	extractor2 "sibyl2/pkg/extractor"
)

var userSrc string
var userLangType string
var userExtractType string

var allowExtractType = []string{"func", "symbol"}

var extractCmd = &cobra.Command{
	Use:    "extract",
	Short:  "test",
	Long:   `test`,
	Hidden: false,
	Run: func(cmd *cobra.Command, args []string) {
		langType := core.LangTypeValueOf(userLangType)
		if langType == core.LangUnknown {
			panic(errors.New("unknown lang type: " + userLangType))
		}

		if !slices.Contains(allowExtractType, userExtractType) {
			panic(errors.New("non-allow extract type: " + userExtractType))
		}

		// call
		runner := &core.Runner{}
		fileUnits, err := runner.File2Units(userSrc, langType)
		if err != nil {
			panic(err)
		}

		extractor := extractor2.GetExtractor(core.LangGo)
		for _, eachFileUnit := range fileUnits {
			if userExtractType == "func" {
				_, err = extractor.ExtractFunctions(eachFileUnit.Units)
			} else {
				_, err = extractor.ExtractSymbols(eachFileUnit.Units)
			}
			if err != nil {
				panic(err)
			}
		}
		// todo: how to collect these results?
	},
}

func init() {
	extractCmd.PersistentFlags().StringVar(&userSrc, "src", ".", "src dir path")

	extractCmd.PersistentFlags().StringVar(&userLangType, "lang", "", "lang type of your source code")
	err := extractCmd.MarkPersistentFlagRequired("lang")
	if err != nil {
		panic(err)
	}

	extractCmd.PersistentFlags().StringVar(&userExtractType, "type", "", "what kind of data you want")
	err = extractCmd.MarkPersistentFlagRequired("type")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(extractCmd)
}
