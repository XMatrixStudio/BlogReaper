package model

import (
	"github.com/boltdb/bolt"
	"github.com/globalsign/mgo/bson"
	"github.com/kataras/iris/core/errors"
	"time"
)

type PublicModel struct {
	*Model
}

type PublicFeed struct {
	URL        string   `bson:"url"`        // 订阅源的URL
	Title      string   `bson:"title"`      // 订阅源的标题
	Subtitle   string   `bson:"subtitle"`   // 订阅源的子标题
	Articles   []string `bson:"articles"`   // 订阅源包括的文章URL
	Star       int64    `bson:"star"`       // 订阅数量
	UpdateDate int64    `bson:"updateDate"` // 更新时间，如果超过12小时就更新，或强制更新
}

type PublicArticle struct {
	URL        string   `bson:"url"`
	FeedURL    string   `bson:"feedUrl"`
	Title      string   `bson:"title"`
	Published  string   `bson:"published"`
	Updated    string   `bson:"updated"`
	Content    string   `bson:"content"`
	Summary    string   `bson:"summary"`
	Categories []string `bson:"categories"`
	Read       int64    `bson:"read"`
}

func (m *PublicModel) AddOrUpdatePublicFeed(url, title, subtitle string, articles []string) error {
	return m.Update(func(b *bolt.Bucket) error {
		fb, err := b.CreateBucketIfNotExists([]byte("feed"))
		if err != nil {
			return err
		}
		bytes := fb.Get([]byte(url))
		if bytes == nil {
			bytes, err = bson.Marshal(&PublicFeed{
				URL:        url,
				Title:      title,
				Subtitle:   subtitle,
				Articles:   articles,
				Star:       0,
				UpdateDate: time.Now().Unix(),
			})
			if err != nil {
				return err
			}
		} else {
			publicFeed := PublicFeed{}
			err = bson.Unmarshal(bytes, &publicFeed)
			if err != nil {
				return err
			}
			publicFeed.Title = title
			publicFeed.Subtitle = subtitle
			publicFeed.Articles = articles
			publicFeed.UpdateDate = time.Now().Unix()
			bytes, err = bson.Marshal(&publicFeed)
			if err != nil {
				return err
			}
		}
		return fb.Put([]byte(url), bytes)
	})
}

func (m *PublicModel) AddOrUpdatePublicArticles(url string, articles []PublicArticle) (err error) {
	return m.Update(func(b *bolt.Bucket) error {
		ab, err := b.CreateBucketIfNotExists([]byte("article"))
		if err != nil {
			return err
		}
		fb, err := b.CreateBucketIfNotExists([]byte("feed"))
		if err != nil {
			return err
		}
		bytes := fb.Get([]byte(url))
		if bytes == nil {
			for _, v := range articles {
				bytes, err := bson.Marshal(&v)
				if err != nil {
					return err
				}
				err = ab.Put([]byte(v.URL), bytes)
				if err != nil {
					return err
				}
			}
		} else {
			publicFeed := PublicFeed{}
			urlReadNum := make(map[string]int64)
			err = bson.Unmarshal(bytes, &publicFeed)
			if err != nil {
				return err
			}
			for _, v := range publicFeed.Articles {
				bytes = ab.Get([]byte(v))
				if bytes != nil {
					article := PublicArticle{}
					err = bson.Unmarshal(bytes, &article)
					if err != nil {
						return err
					}
					urlReadNum[v] = article.Read
					ab.Delete([]byte(v))
				}
			}
			for _, v := range articles {
				v.Read = urlReadNum[v.URL]
				bytes, err = bson.Marshal(&v)
				if err != nil {
					return err
				}
				err = ab.Put([]byte(v.URL), bytes)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (m *PublicModel) GetPublicFeedByURL(url string) (publicFeed PublicFeed, err error) {
	return publicFeed, m.View(func(b *bolt.Bucket) error {
		fb := b.Bucket([]byte("feed"))
		if fb == nil {
			return errors.New("not_found")
		}
		bytes := fb.Get([]byte(url))
		if bytes == nil {
			return errors.New("not_found")
		}
		return bson.Unmarshal(bytes, &publicFeed)
	})
}

func (m *PublicModel) GetPublicArticleByURL(url string) (article PublicArticle, err error) {
	return article, m.View(func(b *bolt.Bucket) error {
		ab := b.Bucket([]byte("article"))
		if ab == nil {
			return errors.New("not_found")
		}
		bytes := ab.Get([]byte(url))
		if bytes == nil {
			return errors.New("not_found")
		}
		return bson.Unmarshal(bytes, &article)
	})
}
