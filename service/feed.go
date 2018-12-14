package service

import (
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/kataras/iris/core/errors"
	"sort"
)

type FeedService interface {
	GetModel() *model.FeedModel
	AddFeed(userID, id, categoryID string) (feed graphql.Feed, err error)
	GetFeedsByCategoryID(userID, categoryID string) (feeds []graphql.Feed, err error)
	EditFeed(userID, feedID string, title *string, categoryIDs []string) (success bool, err error)
	RemoveFeed(userID, feedID string) (success bool, err error)
	EditArticle(userID, feedID, url string, read, later *bool) (success bool, err error)
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
		for k := range feed.Articles {
			feed.Articles[k].FeedID = v.ID.Hex()
		}
		for _, priv := range v.Articles {
			var flag bool
			for pubk := range feed.Articles {
				if priv.URL == feed.Articles[pubk].URL {
					feed.Articles[pubk].Later = priv.Later
					feed.Articles[pubk].Read = priv.Read
					flag = true
					break
				}
			}
			if flag == false {
				feed.Articles = append(feed.Articles, graphql.Article{
					URL:        priv.URL,
					Title:      priv.Content.Title,
					Published:  priv.Content.Published,
					Updated:    priv.Content.Updated,
					Content:    priv.Content.Content,
					Summary:    priv.Content.Summary,
					Categories: priv.Content.Categories,
					Read:       priv.Read,
					Later:      priv.Later,
					FeedID:     v.ID.Hex(),
				})
			}
		}
		sort.Slice(feed.Articles, func(i, j int) bool {
			return feed.Articles[i].Published >= feed.Articles[j].Published
		})
		var articles []model.Article
		for _, av := range feed.Articles {
			articles = append(articles, model.Article{
				URL:   av.URL,
				Read:  av.Read,
				Later: av.Later,
			})
		}
		err = s.Model.UpdateArticles(userID, v.ID.Hex(), articles)
		if err != nil {
			return nil, err
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
	panic("not implement")
}

func (s *feedService) RemoveFeed(userID, feedID string) (success bool, err error) {
	panic("not implement")
}

func (s *feedService) EditArticle(userID, feedID, url string, read, later *bool) (success bool, err error) {
	article, err := s.Model.GetArticleByURL(userID, feedID, url)
	if err != nil {
		return false, errors.New("invalid_feed_or_url")
	}
	var readBool, laterBool bool
	if read != nil {
		readBool = *read
	} else {
		readBool = article.Read
	}
	if later != nil {
		laterBool = *later
	} else {
		laterBool = article.Later
	}
	publicArticle := model.PublicArticle{}
	if laterBool {
		publicArticle, err = s.Service.Public.GetModel().GetPublicArticleByURL(url)
		if err != nil {
			return false, err
		}
	}
	err = s.Model.EditArticle(userID, feedID, url, readBool, laterBool, publicArticle)
	if err != nil {
		return false, err
	}
	return true, nil
}
