package model

import "github.com/globalsign/mgo/bson"

type FeedModel struct {
	*Model
}

type Feed struct {
	ID         bson.ObjectId   `bson:"id"`         // 订阅源的ID
	UserID     bson.ObjectId   `bson:"userId"`     // 订阅用户的ID
	URL        string          `bson:"url"`        // 订阅源的URL
	Title      string          `bson:"title"`      // 订阅源的标题
	Categories []bson.ObjectId `bson:"categories"` // 订阅源的分类
	Articles   []Article       `bson:"articles"`   // 订阅源包括的文章
}

type Article struct {
	URL   string `bson:"url"`
	Read  bool   `bson:"read"`
	Later bool   `bson:"later"`
}

func (m *FeedModel) GetCategoryByName(userID, name string) (category Category, err error) {
	panic("not implement")
}

func (m *FeedModel) AddCategoryByName(userID, name string) (category Category, err error) {
	panic("not implement")
}

func (m *FeedModel) AddFeed(userID, categoryID string, publicFeed PublicFeed) (feed Feed, err error) {
	panic("not implement")
}

func (m *FeedModel) EditFeed(userID, categoryID, url, title string) (feed Feed, err error) {
	panic("not implement")
}

func (m *FeedModel) EditArticle(userID, categoryID, url, articleURL string, read, later bool) (err error) {
	panic("not implement")
}
