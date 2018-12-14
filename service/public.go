package service

import (
	"encoding/xml"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type PublicService interface {
	GetModel() *model.PublicModel
	GetPublicFeed(url string) (feed graphql.Feed, err error)
}

type publicService struct {
	Model   *model.PublicModel
	Service *Service
}

func NewPublicService(s *Service, m *model.PublicModel) PublicService {
	return &publicService{
		Model:   m,
		Service: s,
	}
}

func (s *publicService) GetModel() *model.PublicModel {
	return s.Model
}

// 从数据库中获取PublicFeed
func (s *publicService) GetPublicFeed(url string) (feed graphql.Feed, err error) {
	publicFeed, err := s.Model.GetPublicFeedByURL(url)
	if err != nil && err.Error() != "not_found" {
		return
	}
	if (err != nil && err.Error() == "not_found") || time.Now().Unix()-publicFeed.UpdateDate > 60*60*12 {
		publicFeed, err = s.UpdatePublicFeed(url)
		if err != nil {
			return
		}
	}
	feed = graphql.Feed{
		ID:       "",
		URL:      publicFeed.URL,
		Title:    publicFeed.Title,
		Subtitle: publicFeed.Subtitle,
		Follow:   int(publicFeed.Follow),
		Articles: []graphql.Article{},
	}
	for _, v := range publicFeed.Articles {
		publicArticle, err := s.Model.GetPublicArticleByURL(v)
		if err != nil {
			return feed, err
		}
		feed.Articles = append(feed.Articles, graphql.Article{
			URL:        publicArticle.URL,
			Title:      publicArticle.Title,
			Published:  publicArticle.Published,
			Updated:    publicArticle.Updated,
			Content:    publicArticle.Content,
			Summary:    publicArticle.Summary,
			Categories: publicArticle.Categories,
			Read:       false,
			Later:      false,
			FeedID:     "",
		})
	}
	return
}

type AtomFeed struct {
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Author   AtomAuthor  `xml:"author"`
	Entries  []AtomEntry `xml:"entry"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

type AtomEntry struct {
	Title      string         `xml:"title"`
	Link       AtomLink       `xml:"link"`
	Published  string         `xml:"published"`
	Updated    string         `xml:"updated"`
	Content    string         `xml:"content"`
	Summary    string         `xml:"summary"`
	Categories []AtomCategory `xml:"category"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type AtomCategory struct {
	Term string `xml:"term,attr"`
}

// 从订阅源拉取数据，更新PublicFeed
func (s *publicService) UpdatePublicFeed(url string) (publicFeed model.PublicFeed, err error) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return publicFeed, errors.New(res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	atomFeed := AtomFeed{}
	err = xml.Unmarshal(bytes, &atomFeed)
	if err != nil {
		return
	}
	var articlesUrl []string
	var articles []model.PublicArticle
	for _, v := range atomFeed.Entries {
		var categories []string
		for _, vc := range v.Categories {
			categories = append(categories, vc.Term)
		}
		articlesUrl = append(articlesUrl, v.Link.Href)
		articles = append(articles, model.PublicArticle{
			URL:        v.Link.Href,
			FeedURL:    url,
			Title:      v.Title,
			Published:  v.Published,
			Updated:    v.Updated,
			Content:    v.Content,
			Summary:    v.Summary,
			Categories: categories,
			Read:       0,
		})
	}
	err = s.Model.AddOrUpdatePublicArticles(url, articles)
	if err != nil {
		return
	}
	err = s.Model.AddOrUpdatePublicFeed(url, atomFeed.Title, atomFeed.Subtitle, articlesUrl)
	if err != nil {
		return
	}
	publicFeed, err = s.Model.GetPublicFeedByURL(url)
	return
}
