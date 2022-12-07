package object

type ExecuteConfig struct {
	Port              int
	DbType            DriverType
	Neo4jUri          string
	Neo4jUserName     string
	Neo4jPassword     string
	UploadWorkerCount int
	UploadQueueSize   int
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
	}
}
