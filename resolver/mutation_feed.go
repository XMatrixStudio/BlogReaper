package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *mutationResolver) AddCategory(ctx context.Context, name string) (*graphql.Category, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	category, err := r.Service.Category.AddCategory(userID, name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *mutationResolver) AddFeed(ctx context.Context, id string, categoryId string) (*graphql.Feed, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	feed, err := r.Service.Feed.AddFeed(userID, id, categoryId)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

func (r *mutationResolver) EditArticle(ctx context.Context, url, feedID string, read *bool, later *bool) (bool, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return false, errors.New("not_login")
	}
	if read == nil && later == nil {
		return false, errors.New("invalid_params")
	}
	success, err := r.Service.Feed.EditArticle(userID, feedID, url, read, later)
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *mutationResolver) EditCategory(ctx context.Context, id string, name string) (bool, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return false, errors.New("not_login")
	}
	success, err := r.Service.Category.EditCategory(userID, id, name)
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *mutationResolver) EditFeed(ctx context.Context, id string, title *string, categoryId []string) (bool, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return false, errors.New("not_login")
	}
	_, err := r.Service.Feed.EditFeed(userID, id, title, categoryId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RemoveCategory(ctx context.Context, id string) (bool, error) {
	panic("not implemented")
}

func (r *mutationResolver) RemoveFeed(ctx context.Context, id string) (bool, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return false, errors.New("not_login")
	}
	return r.Service.Feed.RemoveFeed(userID, id)
}
