package binding

import (
	"context"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tikv/client-go/v2/txnkv"
)

type tikvDriver struct {
	client    *txnkv.Client
	addresses []string
}

func initTikvDriver(config object.ExecuteConfig) Driver {
	addresses := strings.Split(config.TikvAddrs, ",")
	return &tikvDriver{
		addresses: addresses,
	}
}

func (t *tikvDriver) GetType() object.DriverType {
	return object.DriverTypeTikv
}

func (t *tikvDriver) InitDriver(_ context.Context) error {
	client, err := txnkv.NewClient(t.addresses)
	if err != nil {
		return err
	}
	t.client = client
	return nil
}

func (t *tikvDriver) DeferDriver() error {
	if err := t.client.Close(); err != nil {
		return err
	}
	t.client = nil
	return nil
}
