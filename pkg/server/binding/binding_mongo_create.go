package binding

import (
	"context"
	"errors"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *mongoDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	// create list of documents
	docs := make([]interface{}, 0, len(f.Units))
	for _, eachFunc := range f.Units {
		doc := bson.M{
			mongoKeyRepo:          wc.RepoId,
			mongoKeyRev:           wc.RevHash,
			mongoKeyPath:          f.Path,
			mongoKeyFuncSignature: eachFunc.GetSignature(),
			mongoKeyFunc:          eachFunc,
		}
		docs = append(docs, doc)
	}

	// insert documents using WriteModel
	models := make([]mongo.WriteModel, 0, len(docs))
	for _, doc := range docs {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	_, err := collection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDriver) CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error {
	// not yet
	return errors.New("not implemented")
}

func (d *mongoDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContextSlim, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	// create document
	doc := bson.M{
		mongoKeyRepo:          wc.RepoId,
		mongoKeyRev:           wc.RevHash,
		mongoKeyPath:          f.Path,
		mongoKeyFuncCtx:       f,
		mongoKeyFuncSignature: f.GetSignature(),
	}

	// insert document
	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)

	// create list of documents
	docs := make([]interface{}, 0, len(c.Units))
	for _, eachClazz := range c.Units {
		doc := bson.M{
			mongoKeyRepo:           wc.RepoId,
			mongoKeyRev:            wc.RevHash,
			mongoKeyPath:           c.Path,
			mongoKeyClazzSignature: eachClazz.GetSignature(),
			mongoKeyClazz:          eachClazz,
		}
		docs = append(docs, doc)
	}

	// insert documents using WriteModel
	models := make([]mongo.WriteModel, 0, len(docs))
	for _, doc := range docs {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	_, err := collection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	// no need
	return nil
}
