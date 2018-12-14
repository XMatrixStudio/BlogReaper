package service

import (
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/kataras/iris/core/errors"
)

type FeedService interface {
	GetModel() *model.FeedModel
	AddFeed(userID, id, categoryID string) (feed graphql.Feed, err error)
	GetFeedsByCategoryID(userID, categoryID string) (feeds []graphql.Feed, err error)
	EditFeed(userID, feedID string, title *string, categoryIDs []string) (success bool, err error)
	RemoveFeed(userID, feedID string) (success bool, err error)
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

func (s *feedService) GetModel() *model.FeedModel {
	return s.Model
}

func (s *feedService) AddFeed(userID, id, categoryID string) (feed graphql.Feed, err error) {
	_, err = s.Service.Category.GetModel().GetCategoryById(userID, categoryID)
	if err != nil {
		return feed, errors.New("invalid_category")
	}
	feed, err = s.Service.Public.GetPublicFeedByID(id)
	if err != nil {
		return feed, errors.New("invalid_id")
	}
	var articlesUrl []string
	for _, v := range feed.Articles {
		articlesUrl = append(articlesUrl, v.URL)
	}
	_, err = s.Model.GetFeedByPublicID(userID, id)
	if err == nil || err.Error() != "not_found" {
		return feed, errors.New("repeat_feed")
	}
	privateFeed, err := s.Model.AddFeed(userID, id, feed.URL, feed.Title, categoryID, articlesUrl)
	if err != nil {
		return
	}
	feed.ID = privateFeed.ID.Hex()
	feed.Title = privateFeed.Title
	err = s.Service.Public.GetModel().IncreasePublicFeedFollow(id)
	return
}

func (s *feedService) GetFeedsByCategoryID(userID, categoryID string) (feeds []graphql.Feed, err error) {
	privateFeeds, err := s.Model.GetFeedsByCategoryID(userID, categoryID)
	if err != nil {
		return
	}
	for _, v := range privateFeeds {
		feed, err := s.Service.Public.GetPublicFeedByID(v.PublicID.Hex())
		if err != nil {
			return feeds, err
		}
		feeds = append(feeds, graphql.Feed{
			ID:       v.ID.Hex(),
			PublicID: feed.PublicID,
			URL:      v.URL,
			Title:    v.Title,
			Subtitle: feed.Subtitle,
			Follow:   feed.Follow,
			Articles: feed.Articles,
		})
	}
	return
}

// 参数为nil表示不修改
func (s *feedService) EditFeed(userID, feedID string, title *string, categoryIDs []string) (success bool, err error) {
	feed, err := s.Model.GetFeedByID(userID, feedID)
	if err != nil {
		return false, err
	}
	if title == nil {
		title = &feed.Title
	}
	_, err = s.Model.EditFeed(userID, feedID, *title, categoryIDs)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *feedService) RemoveFeed(userID, feedID string) (success bool, err error) {
	panic("not implement")
}
