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

func (r *queryResolver) PopularArticles(ctx context.Context, page int, numPerPage int) ([]graphql.Article, error) {
	panic("not implemented")
}

func (r *queryResolver) PopularFeeds(ctx context.Context, page int, numPerPage int) ([]graphql.Feed, error) {
	panic("not implemented")
}
