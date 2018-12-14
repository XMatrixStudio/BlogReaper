package model

import (
	"github.com/boltdb/bolt"
	"github.com/globalsign/mgo/bson"
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
	Title      string   `bson:"title"`
	Published  string   `bson:"published"`
	Updated    string   `bson:"updated"`
	Content    string   `bson:"content"`
	Summary    string   `bson:"summary"`
	Categories []string `bson:"categories"`
	Read       int64    `bson:"read"`
}


// AddPublicFeed 添加公共源
func (m *PublicModel) AddPublicFeed(url, title, subtitle string, articles []string) error {
	return m.Update(func(b *bolt.Bucket) error {
		bytes, err := bson.Marshal(&PublicFeed{
			URL: url,
			Title: title,
			Subtitle: subtitle,
			Articles: articles,
			Star: 0,
			UpdateDate: int64(time.Now().Second()),
		})
		if err != nil {
			return err
		}
		b.Put([]byte(url), bytes)
		return nil
	})
}