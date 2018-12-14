package service

import (
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/kataras/iris/core/errors"
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
	_, err = s.Service.Category.GetModel().GetCategoryById(userID, categoryID)
	if err != nil {
		return feed, errors.New("invalid_id")
	}
	feed, err = s.Service.Public.GetPublicFeed(url)
	if err != nil {
		return feed, errors.New("invalid_url")
	}
	var articlesUrl []string
	for _, v := range feed.Articles {
		articlesUrl = append(articlesUrl, v.URL)
	}
	_, err = s.Model.GetFeedByURL(userID, url)
	if err == nil {
		return feed, errors.New("repeat_url")
	}
	privateFeed, err := s.Model.AddFeed(userID, url, feed.Title, categoryID, articlesUrl)
	if err != nil {
		return
	}
	feed.ID = privateFeed.ID.Hex()
	feed.Title = privateFeed.Title
	// TODO add read number++
	return
}
