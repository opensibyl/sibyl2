package binding

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *mongoDriver) ReadRepos(ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	cur, err := collection.Distinct(ctx, mongoKeyRepo, nil)
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
		doc := &object.FunctionWithSignature{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		functions = append(functions, doc)
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
		doc := &object.FunctionWithSignature{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		if doc.Span.ContainAnyLine(lines...) {
			functions = append(functions, doc)
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
	for k, verify := range rule {
		verifyStr := runtime.FuncForPC(reflect.ValueOf(verify).Pointer()).Name()
		filter["metadata."+k] = bson.M{"$where": fmt.Sprintf("function() { return (%s)(this.metadata.%s); }", verifyStr, k)}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var searchResult []*object.FunctionWithSignature
	for cursor.Next(ctx) {
		doc := &object.FunctionWithSignature{}
		err := cursor.Decode(doc)
		if err != nil {
			return nil, err
		}
		searchResult = append(searchResult, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return searchResult, nil
}

func (d *mongoDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFunc)

	filter := bson.M{
		mongoKeyRepo: wc.RepoId,
		mongoKeyRev:  wc.RevHash,
		mongoKeyFuncSignature: bson.M{
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
		doc := &object.FunctionWithSignature{}
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
		mongoKeyRepo:          wc.RepoId,
		mongoKeyRev:           wc.RevHash,
		mongoKeyFuncSignature: signature,
	}

	doc := &object.FunctionWithSignature{}
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (d *mongoDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag sibyl2.FuncTag, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
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
		doc := &sibyl2.ClazzWithPath{}
		err := cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		classes = append(classes, doc)
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
		doc := &sibyl2.ClazzWithPath{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		if doc.Span.ContainAnyLine(lines...) {
			classes = append(classes, doc)
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
		doc := &sibyl2.FunctionContextSlim{}
		err := cur.Decode(doc)
		if err != nil {
			return nil, err
		}

		if doc.Span.ContainAnyLine(lines...) {
			functionContexts = append(functionContexts, doc)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return functionContexts, nil
}

func (d *mongoDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
	return nil, errors.New("implement me")
}

func (d *mongoDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContextSlim, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection(mongoCollectionFuncCtx)

	filter := bson.M{
		mongoKeyRepo:          wc.RepoId,
		mongoKeyRev:           wc.RevHash,
		mongoKeyFuncSignature: signature,
	}

	// bad design from chatgpt, haha.
	doc := &struct {
		Context *sibyl2.FunctionContextSlim `bson:"funcctx"`
	}{}
	err := collection.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return doc.Context, nil
}
