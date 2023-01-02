package service

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/binding"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/opensibyl/sibyl2/pkg/server/queue"
)

func InitService(_ object.ExecuteConfig, ctx context.Context, driver binding.Driver, q queue.Queue) {
	sharedContext = ctx
	sharedDriver = driver
	sharedQueue = q
}
