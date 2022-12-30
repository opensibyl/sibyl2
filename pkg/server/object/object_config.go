package object

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

type ServerType = string

const (
	StAll     ServerType = "ALL"
	StWorker  ServerType = "WORKER"
	StGateway ServerType = "GATEWAY"
)

type ExecuteConfig struct {
	// server
	Port int        `mapstructure:"port"`
	Mode ServerType `mapstructure:"mode"`

	// binding
	DbType        DriverType `mapstructure:"dbType"`
	Neo4jUri      string     `mapstructure:"neo4JUri"`
	Neo4jUserName string     `mapstructure:"neo4JUserName"`
	Neo4jPassword string     `mapstructure:"neo4JPassword"`
	BadgerPath    string     `mapstructure:"badgerPath"`
	TikvAddrs     string     `mapstructure:"tikvAddrs"`

	// worker
	WorkerCount     int `mapstructure:"workerCount"`
	WorkerQueueSize int `mapstructure:"workerQueueSize"`

	// queue
	QueueType                 QueueType `mapstructure:"queueType"`
	KafkaAddrs                string    `mapstructure:"kafkaAddrs"`
	KafkaFuncTopic            string    `mapstructure:"kafkaFuncTopic"`
	KafkaFuncConsumerGroup    string    `mapstructure:"kafkaFuncConsumerGroup"`
	KafkaFuncCtxTopic         string    `mapstructure:"kafkaFuncCtxTopic"`
	KafkaFuncCtxConsumerGroup string    `mapstructure:"kafkaFuncCtxConsumerGroup"`
}

func DefaultExecuteConfig() ExecuteConfig {
	return ExecuteConfig{
		9876,
		StAll,
		DtBadger,
		"bolt://localhost:7687",
		"neo4j",
		"neo4j",
		"./.sibyl2Storage",
		"127.0.0.1:2379",
		64,
		// each message = 4k, takes nearly 2gb mem
		512_000,
		QueueTypeMemory,
		"10.177.65.230:9092",
		"sibyl-upload-func",
		"sibyl-consumer-func",
		"sibyl-upload-funcctx",
		"sibyl-consumer-funcctx",
	}
}

func (config *ExecuteConfig) ToJson() (string, error) {
	m, err := config.ToMap()
	if err != nil {
		return "", err
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}

func (config *ExecuteConfig) ToMap() (map[string]any, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(config, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
