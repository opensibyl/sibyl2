package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configPath = "."
	configFile = "sibyl-server-config.json"
	configType = "json"
)

func NewServerCmd() *cobra.Command {
	var writeConfigFile bool
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "sibyl server cmd",
		Long:  `sibyl server cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			config := object.DefaultExecuteConfig()

			// read from config
			viper.AddConfigPath(configPath)
			viper.SetConfigFile(configFile)
			viper.SetConfigType(configType)

			core.Log.Infof("trying to read config from: %s/%s", configPath, configFile)
			err := viper.ReadInConfig()
			if err != nil {
				core.Log.Warnf("no config file found, use default")
			} else {
				core.Log.Infof("found config file")
				err = viper.Unmarshal(&config)

				if err != nil {
					core.Log.Errorf("failed to parse config")
					panic(err)
				}
			}

			// really sicks about reading config from env with viper
			// https://github.com/spf13/viper/issues/584
			// create some shortcuts for mongo ...
			if dbType, ok := os.LookupEnv("SIBYL2_BINDING_DBTYPE"); ok {
				config.DbType = dbType
			}
			if mongoDbName, ok := os.LookupEnv("SIBYL2_BINDING_MONGO_DB_NAME"); ok {
				config.MongoDbName = mongoDbName
			}
			if mongoUri, ok := os.LookupEnv("SIBYL2_BINDING_MONGO_URI"); ok {
				config.MongoURI = mongoUri
			}

			if writeConfigFile {
				// save it back
				// viper has a bug here ... unmarshal is case-insensitively
				// https://github.com/spf13/viper/issues/1014
				// so does json:
				// https://stackoverflow.com/questions/49006073/json-unmarshal-struct-case-sensitively
				usedConfigMap, err := config.ToMap()
				if err != nil {
					panic(err)
				}
				err = viper.MergeConfigMap(usedConfigMap)
				if err != nil {
					panic(err)
				}
				err = viper.WriteConfigAs(viper.ConfigFileUsed())
				if err != nil {
					core.Log.Warnf("failed to write config back")
				}

				// will not run
				return
			}

			ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			err = server.Execute(config, ctx)
			if err != nil {
				panic(err)
			}
		},
	}
	serverCmd.PersistentFlags().BoolVar(&writeConfigFile, "config", false, "write config file back if enabled")

	return serverCmd
}
