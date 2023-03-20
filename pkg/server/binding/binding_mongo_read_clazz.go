package binding

import (
	"context"
	"errors"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeyPath: path,
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var classes []*sibyl2.ClazzWithPath
	for cur.Next(ctx) {
		doc := &MongoFactClazz{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		classes = append(classes, doc.ToClazzWithPath())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (d *mongoDriver) ReadClassesWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionClazz)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeyPath: path,
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var classes []*sibyl2.ClazzWithPath
	for cur.Next(ctx) {
		doc := &MongoFactClazz{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		c := doc.ToClazzWithPath()
		if c.Span.ContainAnyLine(lines...) {
			classes = append(classes, c)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (d *mongoDriver) ReadClassesWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	return nil, errors.New("implement me")
}
