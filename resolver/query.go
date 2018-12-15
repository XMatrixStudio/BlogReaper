package resolver

import (
	"context"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/kataras/iris/core/errors"
)

func (r *queryResolver) User(ctx context.Context) (*graphql.User, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	user, err := r.Service.User.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *queryResolver) Feeds(ctx context.Context, keyword string) ([]graphql.Feed, error) {
	if keyword == "" {
		return nil, errors.New("empty_keyword")
	}
	feeds, err := r.Service.Public.GetPublicFeedByKeyword(keyword)
	return feeds, err
}

func (r *queryResolver) PopularArticles(ctx context.Context, page int, numPerPage int) ([]graphql.Article, error) {
	if page <= 0 || numPerPage <= 0 {
		return nil, errors.New("invalid_params")
	}
	feeds, err := r.Service.Public.GetPopularPublicArticles(page, numPerPage)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (r *queryResolver) PopularFeeds(ctx context.Context, page int, numPerPage int) ([]graphql.Feed, error) {
	if page <= 0 || numPerPage <= 0 {
		return nil, errors.New("invalid_params")
	}
	feeds, err := r.Service.Public.GetPopularPublicFeeds(page, numPerPage)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}
