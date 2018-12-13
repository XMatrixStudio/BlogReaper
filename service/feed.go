package service

import (
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
)

type FeedService interface {
	AddFeed(userID, url, categoryID string) (feed graphql.Feed, err error)
}

type feedService struct {
	Model   *model.FeedModel
	Service *Service
}

func NewFeedService(s *Service, m *model.FeedModel) FeedService {
	return &feedService{
		Model:   m,
		Service: s,
	}
}

func (s *feedService) AddFeed(userID, url, categoryID string) (feed graphql.Feed, err error) {
	panic("not implement")
}
