package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

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
	panic("not implement")
}
