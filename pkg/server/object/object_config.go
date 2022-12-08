package object

type ExecuteConfig struct {
	Port int `json:"port"`
	// binding
	DbType        DriverType `json:"dbType"`
	Neo4jUri      string     `json:"neo4JUri"`
	Neo4jUserName string     `json:"neo4JUserName"`
	Neo4jPassword string     `json:"neo4JPassword"`
	// worker
	WorkerCount     int `json:"workerCount"`
	WorkerQueueSize int `json:"workerQueueSize"`
	// queue
	QueueType         QueueType `json:"queueType"`
	KafkaAddrs        string    `json:"kafkaAddrs"`
	KafkaFuncTopic    string    `json:"kafkaFuncTopic"`
	KafkaFuncCtxTopic string    `json:"kafkaFuncCtxTopic"`
}

func DefaultExecuteConfig() ExecuteConfig {
	return ExecuteConfig{
		9876,
		DtInMemory,
		"bolt://localhost:7687",
		"neo4j",
		"neo4j",
		64,
		1024,
		QueueTypeMemory,
		"localhost:9092,localhost:9093,localhost:9094",
		"sibyl-upload-func",
		"sibyl-upload-funcctx",
	}
}
