package queue

import (
	"context"

	"github.com/williamfzc/sibyl2/pkg/server/object"
)

type Queue interface {
	SubmitFunc(unit *object.FunctionUploadUnit)
	SubmitFuncCtx(unit *object.FunctionContextUploadUnit)
	WatchFunc(chan<- *object.FunctionUploadUnit)
	WatchFuncCtx(chan<- *object.FunctionContextUploadUnit)
}

func InitQueue(_ object.ExecuteConfig, _ context.Context) Queue {
	return &MemoryQueue{}
}
