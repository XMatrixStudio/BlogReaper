package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *queryResolver) Feeds(ctx context.Context, keyword string) ([]graphql.Feed, error) {
	if keyword == "" {
		return nil, errors.New("empty_keyword")
	}
	feeds, err := r.Service.Public.GetPublicFeedByKeyword(keyword)
	return feeds, err
}
