package binding

import (
	"context"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) DeleteWorkspace(wc *object.WorkspaceConfig, ctx context.Context) error {
	collection := d.client.Database(d.config.MongoDbName).Collection("clazz_files")

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
	}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
