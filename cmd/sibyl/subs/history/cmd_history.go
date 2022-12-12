package history

import (
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/pkg/core"
)

/*

source code history visualization
inspired by https://github.com/acaudwell/Gource

why:

- logical layer info
- better extensible, with golang
- no heavy render dependencies

how it works:

- pick the first version, create base graph
- for each rev:
	- copy a new one from base graph
	- diff with the previous one, get the diff files
	- create graph for these files
	- update these files on graph
	- save this graph
- convert these graphs to video
*/

var sourceDir string
var outputFile string
var full bool

func NewHistoryCmd() *cobra.Command {
	historyCmd := &cobra.Command{
		Use:    "history",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			err := handle(sourceDir, outputFile, full)
			if err != nil {
				core.Log.Errorf("error when handle: %v", err)
				panic(err)
			}
		},
	}

	historyCmd.PersistentFlags().StringVar(&sourceDir, "src", ".", "src dir path")
	historyCmd.PersistentFlags().StringVar(&outputFile, "output", "./output.html", "output.html")
	historyCmd.PersistentFlags().BoolVar(&full, "full", false, "generate full graph all the time")

	return historyCmd
}
