package server

import (
	"github.com/spf13/cobra"
	"github.com/williamfzc/sibyl2/pkg/server"
	"github.com/williamfzc/sibyl2/pkg/server/binding"
)

var backendUri string
var serverUser string
var serverPwd string

func NewServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "sibyl server cmd",
		Long:  `sibyl server cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			config := server.DefaultExecuteConfig()
			if backendUri != "" {
				config.DbType = binding.DtNeo4j
				config.Neo4jUri = backendUri
			}
			if serverUser != "" {
				config.Neo4jUserName = serverUser
			}
			if serverPwd != "" {
				config.Neo4jPassword = serverPwd
			}

			server.Execute(config)
		},
	}
	serverCmd.PersistentFlags().StringVar(&backendUri, "uri", "", "neo4j backend url")
	serverCmd.PersistentFlags().StringVar(&serverUser, "user", "", "neo4j user")
	serverCmd.PersistentFlags().StringVar(&serverPwd, "pwd", "", "neo4j password")
	return serverCmd
}
