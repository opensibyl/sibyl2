package binding

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.FunctionServiceDTO, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

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

	var functions []*object.FunctionServiceDTO
	for cur.Next(ctx) {
		doc := &MongoFactFunc{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		functions = append(functions, doc.ToFuncWithSignature())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return functions, nil
}

func (d *mongoDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionServiceDTO, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

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

	var functions []*object.FunctionServiceDTO
	for cur.Next(ctx) {
		doc := &MongoFactFunc{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		f := doc.ToFuncWithSignature()
		if f.Span.ContainAnyLine(lines...) {
			functions = append(functions, f)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return functions, nil
}

func (d *mongoDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionServiceDTO, error) {
	if len(rule) == 0 {
		return nil, errors.New("rule is empty")
	}

	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	final := make([]*object.FunctionServiceDTO, 0)
	for cur.Next(ctx) {
		val := &MongoFactFunc{}
		if err := cur.Decode(val); err != nil {
			return nil, err
		}
		fws := val.ToFuncWithSignature()

		d, err := json.Marshal(fws)
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
			final = append(final, fws)
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return final, nil
}

func (d *mongoDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeySignature: bson.M{
			"$regex": regex,
		},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var signatures []string
	for cur.Next(ctx) {
		doc := &MongoFactFunc{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		signatures = append(signatures, doc.Signature)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return signatures, nil
}

func (d *mongoDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionServiceDTO, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo:      wc.RepoId,
		mongoKeyRev:       wc.RevHash,
		mongoKeySignature: signature,
	}

	doc := &MongoFactFunc{}
	err := collection.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		return nil, err
	}

	return doc.ToFuncWithSignature(), nil
}

func (d *mongoDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag object.FuncTag, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeyTag:  bson.M{"$in": []string{tag}},
	}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var functions []string
	for cursor.Next(ctx) {
		var function bson.M
		err = cursor.Decode(&function)
		if err != nil {
			return nil, err
		}

		functionSignature, ok := function[mongoKeySignature].(string)
		if !ok {
			return nil, errors.New("function Signature not found")
		}

		functions = append(functions, functionSignature)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return functions, nil
}
