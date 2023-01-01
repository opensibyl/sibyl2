package queue

import (
	"context"
	"strings"
	"time"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/segmentio/kafka-go"
)

type KafkaQueue struct {
	funcPushList    []chan<- *object.FunctionUploadUnit
	funcCtxPushList []chan<- *object.FunctionContextUploadUnit
	clazzPushList   []chan<- *object.ClazzUploadUnit

	ctx                context.Context
	kafkaWriter        *kafka.Writer
	kafkaFuncReader    *kafka.Reader
	kafkaFuncCtxReader *kafka.Reader
	kafkaClazzReader   *kafka.Reader
	funcTopic          string
	funcCtxTopic       string
	clazzTopic         string
}

func (k *KafkaQueue) GetType() object.QueueType {
	return object.QueueTypeKafka
}

func (k *KafkaQueue) Defer() error {
	return k.kafkaWriter.Close()
}

func (k *KafkaQueue) SubmitFunc(unit *object.FunctionUploadUnit) error {
	v, err := object.SerializeUploadUnit(unit)
	if err != nil {
		core.Log.Errorf("error when serialize upload unit: %v", err)
		return err
	}

	// why this `write` a little slow ...
	err = k.kafkaWriter.WriteMessages(k.ctx, kafka.Message{
		Topic: k.funcTopic,
		Value: v,
		Time:  time.Time{},
	})
	if err != nil {
		core.Log.Errorf("error when write kafka msg: %v", err)
		return err
	}
	return nil
}

func (k *KafkaQueue) SubmitFuncCtx(unit *object.FunctionContextUploadUnit) error {
	v, err := object.SerializeUploadUnit(unit)
	if err != nil {
		core.Log.Errorf("error when serialize upload unit: %v", err)
		return err
	}

	err = k.kafkaWriter.WriteMessages(k.ctx, kafka.Message{
		Topic: k.funcCtxTopic,
		Value: v,
		Time:  time.Time{},
	})
	if err != nil {
		core.Log.Errorf("error when write kafka msg: %v", err)
		return err
	}
	return nil
}

func (k *KafkaQueue) SubmitClazz(unit *object.ClazzUploadUnit) (err error) {
	v, err := object.SerializeUploadUnit(unit)
	if err != nil {
		core.Log.Errorf("error when serialize upload unit: %v", err)
		return err
	}

	err = k.kafkaWriter.WriteMessages(k.ctx, kafka.Message{
		Topic: k.clazzTopic,
		Value: v,
		Time:  time.Time{},
	})
	if err != nil {
		core.Log.Errorf("error when write kafka msg: %v", err)
		return err
	}
	return nil
}

func (k *KafkaQueue) WatchFunc(units chan<- *object.FunctionUploadUnit) {
	go func() {
		for {
			m, err := k.kafkaFuncReader.ReadMessage(k.ctx)
			core.Log.Debugf("rece new func: %d", m.Offset)
			if err != nil {
				core.Log.Errorf("kafka read failed: %v", err)
				break
			}
			unit, err := object.DeserializeFuncUploadUnit(m.Value)
			if err != nil {
				core.Log.Warnf("not a valid func upload object: %v", err)
			}
			units <- unit
		}
	}()
}

func (k *KafkaQueue) WatchFuncCtx(units chan<- *object.FunctionContextUploadUnit) {
	go func() {
		for {
			m, err := k.kafkaFuncCtxReader.ReadMessage(k.ctx)
			core.Log.Debugf("rece new func ctx: %d", m.Offset)
			if err != nil {
				core.Log.Errorf("kafka read failed: %v", err)
				break
			}
			unit, err := object.DeserializeFuncCtxUploadUnit(m.Value)
			if err != nil {
				core.Log.Warnf("not a valid func ctx upload object: %v", err)
			}
			units <- unit
		}
	}()
}

func (k *KafkaQueue) WatchClazz(units chan<- *object.ClazzUploadUnit) {
	go func() {
		for {
			m, err := k.kafkaClazzReader.ReadMessage(k.ctx)
			core.Log.Debugf("rece new clazz: %d", m.Offset)
			if err != nil {
				core.Log.Errorf("kafka read failed: %v", err)
				break
			}
			unit, err := object.DeserializeClazzUploadUnit(m.Value)
			if err != nil {
				core.Log.Warnf("not a valid clazz upload object: %v", err)
			}
			units <- unit
		}
	}()
}

func newKafkaQueue(config object.ExecuteConfig, ctx context.Context) *KafkaQueue {
	addr := strings.Split(config.KafkaAddrs, ",")

	// todo writer and reader in different inst
	funcWriter := &kafka.Writer{
		Addr:     kafka.TCP(addr...),
		Balancer: &kafka.LeastBytes{},
	}
	funcReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: addr,
		Topic:   config.KafkaFuncTopic,
		GroupID: config.KafkaFuncConsumerGroup,
	})
	funcCtxReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: addr,
		Topic:   config.KafkaFuncCtxTopic,
		GroupID: config.KafkaFuncCtxConsumerGroup,
	})
	clazzReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: addr,
		Topic:   config.KafkaClazzTopic,
		GroupID: config.KafkaClazzConsumerGroup,
	})

	return &KafkaQueue{
		kafkaWriter:        funcWriter,
		ctx:                ctx,
		funcTopic:          config.KafkaFuncTopic,
		funcCtxTopic:       config.KafkaFuncCtxTopic,
		clazzTopic:         config.KafkaClazzTopic,
		kafkaFuncReader:    funcReader,
		kafkaFuncCtxReader: funcCtxReader,
		kafkaClazzReader:   clazzReader,
	}
}
