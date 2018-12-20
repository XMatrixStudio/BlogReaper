package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *mutationResolver) AddPublicFeedOrNot(ctx context.Context, url string) (*graphql.Feed, error) {
	feed, err := r.Service.Public.GetPublicFeedByURL(url)
	if err != nil {
		return nil, errors.New("invalid_url")
	}
	return &feed, nil
}
