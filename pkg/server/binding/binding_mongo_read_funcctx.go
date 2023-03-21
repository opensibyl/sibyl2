package binding

import (
	"context"
	"encoding/json"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *mongoDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionContextSlim, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

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

	var functionContexts []*object.FunctionContextSlim
	for cur.Next(ctx) {
		doc := &MongoRelFuncCtx{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		f := doc.ToFuncCtx()
		if f.Span.ContainAnyLine(lines...) {
			functionContexts = append(functionContexts, f)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return functionContexts, nil
}

func (d *mongoDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionContextSlim, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	final := make([]*object.FunctionContextSlim, 0)
	for cur.Next(ctx) {
		val := &MongoRelFuncCtx{}
		if err := cur.Decode(val); err != nil {
			return nil, err
		}

		funcctx := val.ToFuncCtx()
		d, err := json.Marshal(funcctx)
		if err != nil {
			return nil, err
		}

		passed := true
		for rk, verify := range rule {
			v := gjson.GetBytes(d, rk)
			if !verify(v.String()) {
				// failed and ignore this item
				passed = false
				break
			}
		}
		// all the rules passed
		if passed {
			final = append(final, funcctx)
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return final, nil
}

func (d *mongoDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionContextSlim, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	filter := bson.M{
		mongoKeyRepo:      wc.RepoId,
		mongoKeyRev:       wc.RevHash,
		mongoKeySignature: signature,
	}

	doc := &MongoRelFuncCtx{}
	err := collection.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return doc.ToFuncCtx(), nil
}
