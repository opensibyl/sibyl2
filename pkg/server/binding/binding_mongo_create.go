package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *mongoDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	if f.IsEmpty() {
		return nil
	}

	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	// create list of documents
	docs := make([]interface{}, 0, len(f.Units))
	for _, eachFunc := range f.Units {
		doc := &MongoFactFunc{
			MongoFactBase: &MongoFactBase{
				RepoId:    wc.RepoId,
				RevHash:   wc.RevHash,
				Path:      f.Path,
				Signature: eachFunc.GetSignature(),
				Tags:      []string{},
			},
			Func: eachFunc,
		}
		docs = append(docs, doc)
	}

	// insert documents using WriteModel
	models := make([]mongo.WriteModel, 0, len(docs))
	for _, doc := range docs {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	_, err := collection.BulkWrite(ctx, models)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	return nil
}

func (d *mongoDriver) CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	filter := bson.M{
		mongoKeyRepo:      wc.RepoId,
		mongoKeyRev:       wc.RevHash,
		mongoKeySignature: signature,
	}

	update := bson.M{
		"$addToSet": bson.M{
			mongoKeyTag: tag,
		},
	}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

func (d *mongoDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *object.FunctionContextSlim, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	// create document
	doc := &MongoRelFuncCtx{
		MongoFactBase: &MongoFactBase{
			RepoId:    wc.RepoId,
			RevHash:   wc.RevHash,
			Path:      f.Path,
			Signature: f.GetSignature(),
			Tags:      []string{},
		},
		FuncCtx: f,
	}

	// insert document
	_, err := collection.InsertOne(ctx, doc)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

func (d *mongoDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error {
	if c.IsEmpty() {
		return nil
	}

	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)

	// create list of documents
	docs := make([]interface{}, 0, len(c.Units))
	for _, eachClazz := range c.Units {
		doc := &MongoFactClazz{
			MongoFactBase: &MongoFactBase{
				RepoId:    wc.RepoId,
				RevHash:   wc.RevHash,
				Path:      c.Path,
				Signature: eachClazz.GetSignature(),
				Tags:      []string{},
			},
			Clazz: eachClazz,
		}
		docs = append(docs, doc)
	}

	// insert documents using WriteModel
	models := make([]mongo.WriteModel, 0, len(docs))
	for _, doc := range docs {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	_, err := collection.BulkWrite(ctx, models)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}

func (d *mongoDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	// no need
	return nil
}
