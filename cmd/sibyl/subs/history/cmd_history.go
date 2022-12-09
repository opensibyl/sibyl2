package history

import (
	"github.com/spf13/cobra"
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

func NewHistoryCmd() *cobra.Command {
	historyCmd := &cobra.Command{
		Use:    "history",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return historyCmd
}
