package queue

import (
	"context"

	"github.com/williamfzc/sibyl2/pkg/server/object"
)

type Queue interface {
	GetType() object.QueueType
	Defer() error
	SubmitFunc(unit *object.FunctionUploadUnit) (err error)
	SubmitFuncCtx(unit *object.FunctionContextUploadUnit) (err error)
	WatchFunc(chan<- *object.FunctionUploadUnit)
	WatchFuncCtx(chan<- *object.FunctionContextUploadUnit)
}

func InitQueue(config object.ExecuteConfig, ctx context.Context) Queue {
	switch config.QueueType {
	case object.QueueTypeMemory:
		return newMemoryQueue()
	case object.QueueTypeKafka:
		return newKafkaQueue(config, ctx)
	default:
		return newMemoryQueue()
	}
}
