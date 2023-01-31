package frontend

import (
	"github.com/opensibyl/sibyl2/pkg/server"
	"github.com/spf13/cobra"
)

func NewFrontendCmd() *cobra.Command {
	var port int
	var serverCmd = &cobra.Command{
		Use:   "frontend",
		Short: "sibyl frontend cmd",
		Long:  `sibyl frontend cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			err := server.ExecuteFrontend(port, cmd.Context())
			if err != nil {
				panic(err)
			}
		},
	}

	serverCmd.PersistentFlags().IntVar(&port, "port", 3000, "frontend port")
	return serverCmd
}
