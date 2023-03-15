package binding

import (
	"context"
	"errors"

	"github.com/opensibyl/sibyl2"
	"github.com/opensibyl/sibyl2/pkg/server/object"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoDriver) ReadRepos(ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadRevs(repoId string, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadRevInfo(wc *object.WorkspaceConfig, ctx context.Context) (*object.RevInfo, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFiles(wc *object.WorkspaceConfig, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctions(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*object.FunctionWithSignature, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionSignaturesWithRegex(wc *object.WorkspaceConfig, regex string, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*object.FunctionWithSignature, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionsWithTag(wc *object.WorkspaceConfig, tag sibyl2.FuncTag, ctx context.Context) ([]string, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadClasses(wc *object.WorkspaceConfig, path string, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	collection := d.client.Database(d.config.MongoDbName).Collection("clazz_files")

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
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadClassesWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.ClazzWithPath, error) {
	return nil, errors.New("implement me")
}

func (d *mongoDriver) ReadFunctionContextsWithLines(wc *object.WorkspaceConfig, path string, lines []int, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionContextsWithRule(wc *object.WorkspaceConfig, rule Rule, ctx context.Context) ([]*sibyl2.FunctionContextSlim, error) {
	// TODO implement me
	panic("implement me")
}

func (d *mongoDriver) ReadFunctionContextWithSignature(wc *object.WorkspaceConfig, signature string, ctx context.Context) (*sibyl2.FunctionContextSlim, error) {
	// TODO implement me
	panic("implement me")
}
