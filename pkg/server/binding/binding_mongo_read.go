package binding

import (
	"context"
	"errors"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) ReadRepos(ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	cur, err := collection.Distinct(ctx, mongoKeyRepo, bson.D{})
	if err != nil {
		return nil, err
	}

	var repos []string
	for _, repo := range cur {
		if repoStr, ok := repo.(string); ok {
			repos = append(repos, repoStr)
		}
	}

	return repos, nil
}

func (d *mongoDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo: repoId,
	}

	cur, err := collection.Distinct(ctx, mongoKeyRev, filter)
	if err != nil {
		return nil, err
	}

	var revs []string
	for _, rev := range cur {
		if revStr, ok := rev.(string); ok {
			revs = append(revs, revStr)
		}
	}

	return revs, nil
}

func (d *mongoDriver) ReadRevInfo(wc *object.WorkspaceConfig, ctx context.Context) (*object.RevInfo, error) {
	return nil, errors.New("not implemented")
}

func (d *mongoDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
	}

	cur, err := collection.Distinct(ctx, mongoKeyPath, filter)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, file := range cur {
		if fileStr, ok := file.(string); ok {
			files = append(files, fileStr)
		}
	}

	return files, nil
}
