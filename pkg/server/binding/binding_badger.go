package binding

import (
	"context"

	"github.com/dgraph-io/badger/v3"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/server/object"
)

type badgerDriver struct {
	db     *badger.DB
	config object.ExecuteConfig
}

func (d *badgerDriver) InitDriver(_ context.Context) error {
	var dbInst *badger.DB
	var err error

	switch d.config.DbType {
	case object.DriverTypeInMemory:
		dbInst, err = badger.Open(badger.DefaultOptions("").WithInMemory(true))
	case object.DriverTypeBadger:
		core.Log.Infof("trying to open: %s", d.config.BadgerPath)
		dbInst, err = badger.Open(badger.DefaultOptions(d.config.BadgerPath))
	default:
		core.Log.Errorf("db type %v invalid", d.config.DbType)
	}
	if err != nil {
		return err
	}
	d.db = dbInst

	return nil
}

func (d *badgerDriver) DeferDriver() error {
	return d.db.Close()
}

func (d *badgerDriver) GetType() object.DriverType {
	return object.DriverTypeBadger
}

func initBadgerDriver(config object.ExecuteConfig) Driver {
	return &badgerDriver{nil, config}
}
