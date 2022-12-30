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

// default neo4j db may be very slow in I/O
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

	q.WatchFunc(funcUnitQueue)
	q.WatchFuncCtx(funcCtxUnitQueue)

	initWorkers(context, driver)
}

func GetFuncQueueTodoCount() int {
	return len(funcUnitQueue)
}

func GetFuncCtxQueueTodoCount() int {
	return len(funcCtxUnitQueue)
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
				err := driver.CreateFuncContext(result.WorkspaceConfig, each, ctx)
				if err != nil {
					// deadlock easily happen in neo4j when creating complex edges
					// append to the queue
					// should replace with dead message queue
					core.Log.Warnf("err when create ctx for: %v", each.GetSignature())
					funcCtxUnitQueue <- result
				}
			}

		case <-ctx.Done():
			return
		}
	}
}
