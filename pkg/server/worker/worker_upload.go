package worker

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/binding"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/sibyl2/pkg/server/queue"
)

var funcUnitQueue chan *object.FunctionUploadUnit
var funcCtxUnitQueue chan *object.FunctionContextUploadUnit
var clazzUnitQueue chan *object.ClazzUploadUnit

// worker count, db connections count
var workerCount int

// tiny mq, will block request when it is full
// so make it large enough in your scenario
// or add a real mq
// for example,
// each build for repo which contains 3000 files = 3000 jobs in seconds
var workerQueueSize int

func InitWorker(config object.ExecuteConfig, context context.Context, driver binding.Driver, q queue.Queue) {
	workerCount = config.WorkerCount
	workerQueueSize = config.WorkerQueueSize

	funcUnitQueue = make(chan *object.FunctionUploadUnit, workerQueueSize)
	funcCtxUnitQueue = make(chan *object.FunctionContextUploadUnit, workerQueueSize)
	clazzUnitQueue = make(chan *object.ClazzUploadUnit, workerQueueSize)

	q.WatchFunc(funcUnitQueue)
	q.WatchFuncCtx(funcCtxUnitQueue)
	q.WatchClazz(clazzUnitQueue)

	initWorkers(context, driver)
}

func GetFuncQueueTodoCount() int {
	return len(funcUnitQueue)
}

func GetFuncCtxQueueTodoCount() int {
	return len(funcCtxUnitQueue)
}

func GetClazzQueueTodoCount() int {
	return len(clazzUnitQueue)
}

func initWorkers(ctx context.Context, driver binding.Driver) {
	for i := 0; i < workerCount; i++ {
		go func() {
			startWorker(ctx, driver)
		}()
	}
}

func startWorker(ctx context.Context, driver binding.Driver) {
	for {
		select {
		case result := <-funcUnitQueue:
			// failure allowed
			// todo: waste 1 txn
			_ = driver.CreateWorkspace(result.WorkspaceConfig, ctx)

			err := driver.CreateFuncFile(result.WorkspaceConfig, result.FunctionResult, ctx)
			if err != nil {
				core.Log.Errorf("error when upload: %v\n", err)
			}

		case result := <-funcCtxUnitQueue:
			// failure allowed
			// todo: waste 1 txn
			_ = driver.CreateWorkspace(result.WorkspaceConfig, ctx)

			for _, each := range result.FunctionContexts {
				err := driver.CreateFuncContext(result.WorkspaceConfig, object.CompressFunctionContext(each), ctx)
				if err != nil {
					// deadlock easily happen in neo4j when creating complex edges
					// append to the queue
					// should replace with dead message queue
					core.Log.Warnf("err when create ctx for: %v, %v", each.GetSignature(), err)
					funcCtxUnitQueue <- result
				}
			}

		case result := <-clazzUnitQueue:
			// failure allowed
			// todo: waste 1 txn
			_ = driver.CreateWorkspace(result.WorkspaceConfig, ctx)

			err := driver.CreateClazzFile(result.WorkspaceConfig, result.ClazzFileResult, ctx)
			if err != nil {
				core.Log.Errorf("error when upload class: %v\n", err)
			}

		case <-ctx.Done():
			return
		}
	}
}
