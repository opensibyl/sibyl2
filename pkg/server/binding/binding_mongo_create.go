package binding

import (
	"context"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/extractor"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) CreateFuncFile(wc *object.WorkspaceConfig, f *extractor.FunctionFileResult, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) CreateFuncTag(wc *object.WorkspaceConfig, signature string, tag string, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) CreateFuncContext(wc *object.WorkspaceConfig, f *sibyl2.FunctionContextSlim, ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) CreateClazzFile(wc *object.WorkspaceConfig, c *extractor.ClazzFileResult, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection("clazz_files")

	for _, eachClazz := range c.Units {
		// create document
		doc := bson.M{
			mongoKeyRepo:           wc.RepoId,
			mongoKeyRev:            wc.RevHash,
			mongoKeyPath:           c.Path,
			mongoKeyClazzSignature: eachClazz.GetSignature(),
			mongoKeyClazz:          eachClazz,
		}

		// insert document
		_, err := collection.InsertOne(ctx, doc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *mongoDriver) CreateWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	// no need
	return nil
}
