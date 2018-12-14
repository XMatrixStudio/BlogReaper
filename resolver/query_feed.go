package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *queryResolver) Articles(ctx context.Context, later *bool, today *bool) ([]graphql.Article, error) {
	panic("not implemented")
}

func (r *queryResolver) Categories(ctx context.Context) ([]graphql.Category, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	categories, err := r.Service.Category.GetCategories(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *queryResolver) Feeds(ctx context.Context) ([]graphql.Feed, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *queryResolver) PopularArticles(ctx context.Context, page int, numPerPage int) ([]graphql.Article, error) {
	panic("not implemented")
}

func (r *queryResolver) PopularFeeds(ctx context.Context, page int, numPerPage int) ([]graphql.Feed, error) {
	panic("not implemented")
}
