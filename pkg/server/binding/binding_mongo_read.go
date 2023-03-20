package binding

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (d *mongoDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.FunctionWithSignature, error) {
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

	var functions []*object.FunctionWithSignature
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

func (d *mongoDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionWithSignature, error) {
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

	var functions []*object.FunctionWithSignature
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

func (d *mongoDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionWithSignature, error) {
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

	final := make([]*object.FunctionWithSignature, 0)
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

func (d *mongoDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionWithSignature, error) {
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

func (d *mongoDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag sibyl2.FuncTag, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)
	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeyTag:  tag,
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

func (d *mongoDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
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

	var functionContexts []*sibyl2.FunctionContextSlim
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

func (d *mongoDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
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

	final := make([]*sibyl2.FunctionContextSlim, 0)
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

func (d *mongoDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContextSlim, error) {
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
