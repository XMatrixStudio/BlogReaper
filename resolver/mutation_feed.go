package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *mutationResolver) AddFeed(ctx context.Context, url string, categoryId *string, categoryName *string) (*graphql.Category, error) {
	if r.Session.Get("id") == nil {
		return nil, errors.New("not_login")
	}
	// TODO
	panic("not implemented")
}

func (r *mutationResolver) EditArticle(ctx context.Context, url string, read *bool, later *bool) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *mutationResolver) EditCategory(ctx context.Context, id string, name string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *mutationResolver) EditFeed(ctx context.Context, url string, title *string, categoryId *string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *mutationResolver) RemoveFeed(ctx context.Context, url string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}
