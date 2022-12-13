package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/spf13/cobra"
)

var diffSrc string
var diffFrom string
var diffTo string
var diffPatch string
var diffRaw string
var diffOutputFile string
var diffThin bool

func NewDiffCommand() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:    "diff",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			var results *ParseResult
			var err error

			if diffFrom != "" && diffTo != "" {
				core.Log.Infof("diff with rev")

				// about why I use cmd rather than some libs
				// because go-git 's patch has some bugs ...
				gitDiffCmd := exec.Command("git", "diff", diffFrom, diffTo)
				gitDiffCmd.Dir = diffSrc
				data, err := gitDiffCmd.CombinedOutput()
				if err != nil {
					core.Log.Errorf("git cmd error: %s", data)
					panic(err)
				}
				results, err = parsePatchRaw(diffSrc, data)
			} else {
				core.Log.Infof("diff with patch")
				if diffRaw != "" {
					results, err = parsePatchRaw(diffSrc, []byte(diffRaw))
				} else {
					results, err = parsePatch(diffSrc, diffPatch)
				}
			}

			if err != nil {
				panic(err)
			}

			var output []byte
			if diffThin {
				output, err = json.MarshalIndent(results.Flatten(), "", "  ")
			} else {
				output, err = json.MarshalIndent(&results, "", "  ")
			}
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
	diffCmd.PersistentFlags().StringVar(&diffFrom, "from", "", "from rev")
	diffCmd.PersistentFlags().StringVar(&diffTo, "to", "", "to rev")

	diffCmd.PersistentFlags().StringVar(&diffPatch, "patch", "", "patch")
	diffCmd.PersistentFlags().StringVar(&diffRaw, "patchRaw", "", "patch raw")

	// for output
	diffCmd.PersistentFlags().StringVar(&diffOutputFile, "output", "", "output json file")
	diffCmd.PersistentFlags().BoolVar(&diffThin, "thin", false, "")
	return diffCmd
}
