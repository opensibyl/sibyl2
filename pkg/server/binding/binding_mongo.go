package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	idxRowStartSuffix = "span.start.row"
	idxRowEndSuffix   = "span.end.row"

	mongoKeyRepo = "repo_id"
	mongoKeyRev  = "rev_hash"
	mongoKeyPath = "path"

	mongoKeyFuncSignature = "func_signature"
	mongoKeyFunc          = "func"
	mongoKeyFuncRowStart  = mongoKeyFunc + "." + idxRowStartSuffix
	mongoKeyFuncRowEnd    = mongoKeyFunc + "." + idxRowEndSuffix

	mongoKeyClazzSignature = "clazz_signature"
	mongoKeyClazz          = "clazz"
	mongoKeyClazzRowStart  = mongoKeyClazz + "." + idxRowStartSuffix
	mongoKeyClazzRowEnd    = mongoKeyClazz + "." + idxRowEndSuffix

	mongoKeyFuncCtx         = "funcctx"
	mongoKeyFuncCtxRowStart = mongoKeyFuncCtx + "." + idxRowStartSuffix
	mongoKeyFuncCtxRowEnd   = mongoKeyFuncCtx + "." + idxRowEndSuffix

	mongoCollectionClazz   = "fact_clazz"
	mongoCollectionFunc    = "fact_func"
	mongoCollectionFuncCtx = "rel_funcctx"
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

	// ensure indexes
	// create unique index on repoId, revHash, path, signature, and span
	funcCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	keys := bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeyFuncSignature, 1},
		{mongoKeyFuncRowStart, 1},
		{mongoKeyFuncRowEnd, 1},
	}
	index := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, err = funcCollection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	// create unique index on repoId, revHash, path, signature, and span
	clazzCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)
	keys = bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeyClazzSignature, 1},
		{mongoKeyClazzRowStart, 1},
		{mongoKeyClazzRowEnd, 1},
	}
	index = mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, err = clazzCollection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	// create unique index on repoId, revHash, path, signature, and span
	funcCtxCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)
	keys = bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeyFuncSignature, 1},
		{mongoKeyFuncCtxRowStart, 1},
		{mongoKeyFuncCtxRowEnd, 1},
	}
	index = mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, err = funcCtxCollection.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

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
