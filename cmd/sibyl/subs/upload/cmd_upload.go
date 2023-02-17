package upload

import (
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewUploadCmd() *cobra.Command {
	var uploadConfigFile string
	var uploadRepoId string
	var uploadRevHash string
	var uploadSrc string
	var uploadLangType []string
	var uploadUrl string
	var uploadWithCtx bool
	var uploadWithClass bool
	var uploadBatchLimit int
	var uploadDryRun bool
	var uploadDepth int

	uploadCmd := &cobra.Command{
		Use:    "upload",
		Short:  "test",
		Long:   `test`,
		Hidden: false,
		Run: func(cmd *cobra.Command, args []string) {
			config := DefaultConfig()
			defaultConf := DefaultConfig()

			// specific config file?
			viper.SetConfigType(configType)
			if uploadConfigFile != "" {
				core.Log.Infof("specific config file: %s", uploadConfigFile)
				viper.SetConfigFile(uploadConfigFile)
			} else {
				// always search in src dir
				viper.AddConfigPath(uploadSrc)
				viper.SetConfigName(configFile)
			}

			err := viper.ReadInConfig()
			if err != nil {
				core.Log.Warnf("no config file found, use default: %v", err)
			} else {
				core.Log.Infof("found config file: %s", viper.ConfigFileUsed())
				err = viper.Unmarshal(config)
				core.Log.Infof("config from file: %v", viper.AllSettings())

				if err != nil {
					core.Log.Errorf("failed to parse config")
					panic(err)
				}
			}

			// read from cmd and overwrite
			// a little ugly ...
			if uploadRepoId != defaultConf.RepoId {
				config.RepoId = uploadRepoId
			}
			if uploadRevHash != defaultConf.RevHash {
				config.RevHash = uploadRevHash
			}
			if uploadSrc != defaultConf.Src {
				config.Src = uploadSrc
			}
			if len(uploadLangType) != 0 {
				config.Lang = uploadLangType
			}
			if uploadUrl != defaultConf.Url {
				config.Url = uploadUrl
			}
			if uploadWithCtx != defaultConf.WithCtx {
				config.WithCtx = uploadWithCtx
			}
			if uploadWithClass != defaultConf.WithClass {
				config.WithClass = uploadWithClass
			}
			if uploadBatchLimit != defaultConf.Batch {
				config.Batch = uploadBatchLimit
			}
			if uploadDryRun != defaultConf.Dry {
				config.Dry = uploadDryRun
			}
			if uploadDepth != defaultConf.Depth {
				config.Depth = uploadDepth
			}

			// execute
			ExecWithConfig(config)

			// save it back
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
		},
	}

	config := DefaultConfig()
	uploadCmd.PersistentFlags().StringVar(&uploadConfigFile, "config", "", "config file path")
	uploadCmd.PersistentFlags().StringVar(&uploadRepoId, "repoId", config.RepoId, "custom repo id")
	uploadCmd.PersistentFlags().StringVar(&uploadRevHash, "revHash", config.RevHash, "custom rev hash")
	uploadCmd.PersistentFlags().StringVar(&uploadSrc, "src", config.Src, "src dir path")
	uploadCmd.PersistentFlags().StringSliceVar(&uploadLangType, "lang", config.Lang, "lang type of your source code")
	uploadCmd.PersistentFlags().StringVar(&uploadUrl, "url", config.Url, "backend url")
	uploadCmd.PersistentFlags().BoolVar(&uploadWithCtx, "withCtx", config.WithCtx, "with func context")
	uploadCmd.PersistentFlags().BoolVar(&uploadWithClass, "withClass", config.WithClass, "with class")
	uploadCmd.PersistentFlags().IntVar(&uploadBatchLimit, "batch", config.Batch, "each batch size")
	uploadCmd.PersistentFlags().BoolVar(&uploadDryRun, "dry", config.Dry, "dry run without upload")
	uploadCmd.PersistentFlags().IntVar(&uploadDepth, "depth", config.Depth, "upload with history")

	return uploadCmd
}
