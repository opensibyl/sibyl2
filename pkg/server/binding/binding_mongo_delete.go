package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	funcCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	clazzCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)
	funcctxCollection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
	}

	_, err := funcCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	_, err = clazzCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	_, err = funcctxCollection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
