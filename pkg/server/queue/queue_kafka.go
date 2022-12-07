package queue

import "github.com/williamfzc/sibyl2/pkg/server/object"

type KafkaQueue struct {
	funcPushList    []chan<- *object.FunctionUploadUnit
	funcCtxPushList []chan<- *object.FunctionContextUploadUnit
}

func (k *KafkaQueue) GetType() object.QueueType {
	return object.QueueTypeKafka
}

func (k *KafkaQueue) SubmitFunc(unit *object.FunctionUploadUnit) {
	//TODO implement me
	panic("implement me")
}

func (k *KafkaQueue) SubmitFuncCtx(unit *object.FunctionContextUploadUnit) {
	//TODO implement me
	panic("implement me")
}

func (k *KafkaQueue) WatchFunc(units chan<- *object.FunctionUploadUnit) {
	//TODO implement me
	panic("implement me")
}

func (k *KafkaQueue) WatchFuncCtx(units chan<- *object.FunctionContextUploadUnit) {
	//TODO implement me
	panic("implement me")
}
