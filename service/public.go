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
	GetPublicFeedByID(id string) (feed graphql.Feed, err error)
	GetPublicFeedByURL(url string) (feed graphql.Feed, err error)
	GetPublicFeedByKeyword(keyword string) (feeds []graphql.Feed, err error)
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
func (s *publicService) GetPublicFeedByID(id string) (feed graphql.Feed, err error) {
	publicFeed, err := s.Model.GetPublicFeedByID(id)
	if err != nil {
		return
	}
	feed = graphql.Feed{
		ID:       "",
		PublicID: publicFeed.ID.Hex(),
		URL:      publicFeed.URL,
		Title:    publicFeed.Title,
		Subtitle: publicFeed.Subtitle,
		Follow:   int(publicFeed.Follow),
		Articles: nil,
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

func (s *publicService) GetPublicFeedByURL(url string) (feed graphql.Feed, err error) {
	publicFeed, err := s.Model.GetPublicFeedByURL(url)
	if err != nil && err.Error() != "not_found" {
		return
	}
	if err != nil && err.Error() == "not_found" {
		publicFeed, err = s.UpdatePublicFeed("", url)
	} else if time.Now().Unix()-publicFeed.UpdateDate > 60*60*12 {
		publicFeed, err = s.UpdatePublicFeed(publicFeed.ID.Hex(), url)
	}
	if err != nil {
		return
	}
	feed = graphql.Feed{
		ID:       "",
		PublicID: publicFeed.ID.Hex(),
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

func (s *publicService) GetPublicFeedByKeyword(keyword string) (feeds []graphql.Feed, err error) {
	publicFeeds, err := s.Model.GetPublicFeedsByKeyword(keyword)
	if err != nil {
		return
	}
	for _, v := range publicFeeds {
		feed := graphql.Feed{
			ID:       "",
			PublicID: v.ID.Hex(),
			URL:      v.URL,
			Title:    v.Title,
			Subtitle: v.Subtitle,
			Follow:   int(v.Follow),
			Articles: []graphql.Article{},
		}
		for _, v := range v.Articles {
			publicArticle, err := s.Model.GetPublicArticleByURL(v)
			if err != nil {
				return feeds, err
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
		feeds = append(feeds, feed)
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

type RSSFeed struct {
	Title string `xml:"title"`
	Description string `xml:"description"`
	Author string `xml:"author"`
	Items []RSSItem `xml:"author"`
}

type RSSItem struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

// 从订阅源拉取数据，更新PublicFeed
func (s *publicService) UpdatePublicFeed(id, url string) (publicFeed model.PublicFeed, err error) {
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
	if id == "" {
		publicFeed, err = s.Model.AddPublicFeed(url, atomFeed.Title, atomFeed.Subtitle, articlesUrl)
	} else {
		publicFeed, err = s.Model.UpdatePublicFeed(id, atomFeed.Title, atomFeed.Subtitle, articlesUrl)
	}
	return
}
