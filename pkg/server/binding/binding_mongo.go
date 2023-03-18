package binding

import (
	"context"

	"github.com/opensibyl/sibyl2"
	object2 "github.com/opensibyl/sibyl2/pkg/extractor/object"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFactBase struct {
	RepoId    string   `bson:"repo_id"`
	RevHash   string   `bson:"rev_hash"`
	Path      string   `bson:"path"`
	Signature string   `bson:"signature"`
	Tag       []string `bson:"tag"`
}

type MongoFactFunc struct {
	*MongoFactBase `bson:",inline"`
	Func           *object2.Function `bson:"func"`
}

func (f *MongoFactFunc) ToFuncWithSignature() *object.FunctionWithSignature {
	return &object.FunctionWithSignature{
		FunctionWithTag: &sibyl2.FunctionWithTag{
			FunctionWithPath: &sibyl2.FunctionWithPath{
				Function: f.Func,
				Path:     f.Path,
			},
			Tags: f.Tag,
		},
		Signature: f.Signature,
	}
}

type MongoFactClazz struct {
	*MongoFactBase `bson:",inline"`
	Clazz          *object2.Clazz `bson:"clazz"`
}

func (c *MongoFactClazz) ToClazzWithPath() *sibyl2.ClazzWithPath {
	return &sibyl2.ClazzWithPath{
		Clazz: c.Clazz,
		Path:  c.Path,
	}
}

type MongoRelFuncCtx struct {
	*MongoFactBase `bson:",inline"`
	FuncCtx        *sibyl2.FunctionContextSlim `bson:"funcctx"`
}

const (
	idxRowStartSuffix = "span.start.row"
	idxRowEndSuffix   = "span.end.row"

	mongoKeyRepo      = "repo_id"
	mongoKeyRev       = "rev_hash"
	mongoKeyPath      = "path"
	mongoKeySignature = "signature"
	mongoKeyTag       = "tag"

	mongoKeyFunc         = "func"
	mongoKeyFuncRowStart = mongoKeyFunc + "." + idxRowStartSuffix
	mongoKeyFuncRowEnd   = mongoKeyFunc + "." + idxRowEndSuffix

	mongoKeyClazz         = "clazz"
	mongoKeyClazzRowStart = mongoKeyClazz + "." + idxRowStartSuffix
	mongoKeyClazzRowEnd   = mongoKeyClazz + "." + idxRowEndSuffix

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
	// create unique index on RepoId, RevHash, Path, Signature, and span
	funcCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	keys := bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeySignature, 1},
		{mongoKeyFuncRowStart, 1},
		{mongoKeyFuncRowEnd, 1},
	}
	index := mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, _ = funcCollection.Indexes().CreateOne(ctx, index)

	// create unique index on RepoId, RevHash, Path, Signature, and span
	clazzCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)
	keys = bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeySignature, 1},
		{mongoKeyClazzRowStart, 1},
		{mongoKeyClazzRowEnd, 1},
	}
	index = mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, _ = clazzCollection.Indexes().CreateOne(ctx, index)

	// create unique index on RepoId, RevHash, Path, Signature, and span
	funcCtxCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)
	keys = bson.D{
		{mongoKeyRepo, 1},
		{mongoKeyRev, 1},
		{mongoKeyPath, 1},
		{mongoKeySignature, 1},
		{mongoKeyFuncCtxRowStart, 1},
		{mongoKeyFuncCtxRowEnd, 1},
	}
	index = mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetUnique(true),
	}
	_, _ = funcCtxCollection.Indexes().CreateOne(ctx, index)

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
