package resolver

import (
	"context"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *mutationResolver) AddFeed(ctx context.Context, url string, categoryId *string, categoryName *string) (*graphql.Category, error) {
	// TODO
	panic("not implemented")
}

func (r *mutationResolver) EditArticle(ctx context.Context, url string, read *bool, later *bool) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) EditCategory(ctx context.Context, id string, name string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) EditFeed(ctx context.Context, url string, title *string, categoryId *string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveFeed(ctx context.Context, url string) (bool, error) {
	panic("not implemented")
}
