package queue

import "github.com/opensibyl/sibyl2/pkg/server/object"

type MemoryQueue struct {
	funcPushList    []chan<- *object.FunctionUploadUnit
	funcCtxPushList []chan<- *object.FunctionContextUploadUnit
	clazzPushList   []chan<- *object.ClazzUploadUnit
}

func (q *MemoryQueue) GetType() object.QueueType {
	return object.QueueTypeMemory
}

func (q *MemoryQueue) Defer() error {
	// do nothing
	return nil
}

func (q *MemoryQueue) WatchFunc(c chan<- *object.FunctionUploadUnit) {
	q.funcPushList = append(q.funcPushList, c)
}

func (q *MemoryQueue) WatchFuncCtx(c chan<- *object.FunctionContextUploadUnit) {
	q.funcCtxPushList = append(q.funcCtxPushList, c)
}

func (q *MemoryQueue) WatchClazz(c chan<- *object.ClazzUploadUnit) {
	q.clazzPushList = append(q.clazzPushList, c)
}

func (q *MemoryQueue) SubmitFunc(unit *object.FunctionUploadUnit) error {
	for _, each := range q.funcPushList {
		each <- unit
	}
	return nil
}

func (q *MemoryQueue) SubmitFuncCtx(unit *object.FunctionContextUploadUnit) error {
	for _, each := range q.funcCtxPushList {
		each <- unit
	}
	return nil
}

func (q *MemoryQueue) SubmitClazz(unit *object.ClazzUploadUnit) (err error) {
	for _, each := range q.clazzPushList {
		each <- unit
	}
	return nil
}

func newMemoryQueue() *MemoryQueue {
	return &MemoryQueue{}
}
