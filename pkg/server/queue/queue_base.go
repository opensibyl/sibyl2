package queue

import (
	"context"

	"github.com/williamfzc/sibyl2/pkg/server/object"
)

type Queue interface {
	GetType() object.QueueType
	SubmitFunc(unit *object.FunctionUploadUnit)
	SubmitFuncCtx(unit *object.FunctionContextUploadUnit)
	WatchFunc(chan<- *object.FunctionUploadUnit)
	WatchFuncCtx(chan<- *object.FunctionContextUploadUnit)
}

func InitQueue(config object.ExecuteConfig, _ context.Context) Queue {
	switch config.QueueType {
	case object.QueueTypeMemory:
		return &MemoryQueue{}
	case object.QueueTypeKafka:
		return &KafkaQueue{}
	default:
		return &MemoryQueue{}
	}
}
