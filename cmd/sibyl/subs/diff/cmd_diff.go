package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/pkg/core"
)

var diffSrc string
var diffPatch string
var diffRaw string
var diffOutputFile string

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
			if diffOutputFile == "" {
				diffOutputFile = fmt.Sprintf("sibyl-diff-%d.json", time.Now().Unix())
			}
			err = os.WriteFile(diffOutputFile, output, 0644)
			if err != nil {
				panic(err)
			}
			core.Log.Infof("file has been saved to: %s", diffOutputFile)
		},
	}
	diffCmd.PersistentFlags().StringVar(&diffSrc, "src", ".", "src dir path")
	diffCmd.PersistentFlags().StringVar(&diffPatch, "patch", "", "patch")
	diffCmd.PersistentFlags().StringVar(&diffRaw, "patchRaw", "", "patch raw")
	diffCmd.PersistentFlags().StringVar(&diffOutputFile, "output", "", "output json file")
	return diffCmd
}
