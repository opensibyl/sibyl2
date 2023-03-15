package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoKeyRepo           = "repo_id"
	mongoKeyRev            = "rev_hash"
	mongoKeyPath           = "path"
	mongoKeyClazzSignature = "clazz_signature"
	mongoKeyClazz          = "clazz"
)

type mongoDriver struct {
	client *mongo.Client
	config object.ExecuteConfig
}

func (d *mongoDriver) InitDriver(ctx context.Context) error {
	clientInst, err := mongo.Connect(
		ctx, options.Client().ApplyURI(d.config.MongoURI))
	if err != nil {
		return err
	}
	d.client = clientInst

	return nil
}

func (d *mongoDriver) DeferDriver() error {
	return d.client.Disconnect(context.Background())
}

func (d *mongoDriver) GetType() object.DriverType {
	return object.DriverTypeMongoDB
}

func initMongoDriver(config object.ExecuteConfig) Driver {
	return &mongoDriver{nil, config}
}
