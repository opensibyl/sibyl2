package queue

import "github.com/williamfzc/sibyl2/pkg/server/object"

type MemoryQueue struct {
	funcPushList    []chan<- *object.FunctionUploadUnit
	funcCtxPushList []chan<- *object.FunctionContextUploadUnit
}

func (q *MemoryQueue) GetType() object.QueueType {
	return object.QueueTypeMemory
}

func (q *MemoryQueue) WatchFunc(c chan<- *object.FunctionUploadUnit) {
	q.funcPushList = append(q.funcPushList, c)
}

func (q *MemoryQueue) WatchFuncCtx(c chan<- *object.FunctionContextUploadUnit) {
	q.funcCtxPushList = append(q.funcCtxPushList, c)
}

func (q *MemoryQueue) SubmitFunc(unit *object.FunctionUploadUnit) {
	for _, each := range q.funcPushList {
		each <- unit
	}
}

func (q *MemoryQueue) SubmitFuncCtx(unit *object.FunctionContextUploadUnit) {
	for _, each := range q.funcCtxPushList {
		each <- unit
	}
}

func newMemoryQueue() *MemoryQueue {
	return &MemoryQueue{}
}
