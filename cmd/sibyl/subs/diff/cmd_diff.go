package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var diffSrc string
var diffPatch string
var diffRaw string
var diffOutput string

func NewDiffCommand() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:    "diff",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			var results *ParseResult
			var err error
			if diffRaw != "" {
				results, err = parsePatchRaw(diffSrc, []byte(diffRaw))
			} else {
				results, err = parsePatch(diffSrc, diffPatch)
			}
			if err != nil {
				panic(err)
			}

			output, err := json.MarshalIndent(&results, "", "  ")
			if err != nil {
				panic(err)
			}
			if diffOutput == "" {
				diffOutput = fmt.Sprintf("sibyl-diff-%d.json", time.Now().Unix())
			}
			err = os.WriteFile(diffOutput, output, 0644)
			if err != nil {
				panic(err)
			}
		},
	}
	diffCmd.PersistentFlags().StringVar(&diffSrc, "src", ".", "src dir path")
	diffCmd.PersistentFlags().StringVar(&diffPatch, "patch", "", "patch")
	diffCmd.PersistentFlags().StringVar(&diffRaw, "patchRaw", "", "patch raw")
	diffCmd.PersistentFlags().StringVar(&diffOutput, "output", "", "output json file")
	return diffCmd
}
