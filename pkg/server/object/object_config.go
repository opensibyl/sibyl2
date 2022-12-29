package object

import "encoding/json"

type ExecuteConfig struct {
	// server
	Port int `json:"port"`

	// binding
	DbType        DriverType `json:"dbType"`
	Neo4jUri      string     `json:"neo4JUri"`
	Neo4jUserName string     `json:"neo4JUserName"`
	Neo4jPassword string     `json:"neo4JPassword"`
	BadgerPath    string     `json:"badgerPath"`

	// worker
	WorkerCount     int `json:"workerCount"`
	WorkerQueueSize int `json:"workerQueueSize"`

	// queue
	QueueType                 QueueType `json:"queueType"`
	KafkaAddrs                string    `json:"kafkaAddrs"`
	KafkaFuncTopic            string    `json:"kafkaFuncTopic"`
	KafkaFuncConsumerGroup    string    `json:"kafkaFuncConsumerGroup"`
	KafkaFuncCtxTopic         string    `json:"kafkaFuncCtxTopic"`
	KafkaFuncCtxConsumerGroup string    `json:"kafkaFuncCtxConsumerGroup"`
}

func DefaultExecuteConfig() ExecuteConfig {
	return ExecuteConfig{
		9876,
		DtBadger,
		"bolt://localhost:7687",
		"neo4j",
		"neo4j",
		"./.sibyl2Storage",
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
	bytes, err := json.Marshal(config)
	if err != nil {
		return "", nil
	}
	return string(bytes), nil
}

func (config *ExecuteConfig) ToMap() (map[string]any, error) {
	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
