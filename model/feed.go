package model

import (
	"errors"
	"github.com/globalsign/mgo/bson"
	"sort"
	"encoding/json"
)

type FeedModel struct {
	*Model
}

type Feed struct {
	UserID     string     `json:"uid" db:"uid"`               // 用户ID
	ID         string     `json:"id" db:"id"`                 // 订阅源的ID
	PublicID   string     `json:"pid" db:"pid"`     // 订阅源的公共ID
	URL        string     `json:"url" db:"url"`               // 订阅源的URL
	Title      string     `json:"title" db:"title"`           // 订阅源的标题
	Categories []string   `json:"categories" db:"categories"` // 订阅源的分类
	Articles   []Article  `json:"articles" db:"articles"`     // 订阅源包括的文章
}

type Article struct {
	URL     string         `json:"articleurl" db:"articleurl"`
	Read    bool           `json:"read" db:"read"`
	Later   bool           `json:"later" db:"later"`
	Content *PublicArticle `json:"publicarticle" db:"publicarticle"`
}

func (m *FeedModel) AddFeed(userID, publicID, url, title, categoryID string, articlesUrl []string) (feed Feed, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(publicID) || !bson.IsObjectIdHex(categoryID) {
		return feed, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`INSERT INTO feed VALUES(?)`)
	if err != nil {
		return feed, err
	}
	var articles []Article
	for _, a := range articlesUrl {
		articles = append(articles, Article{
			URL:     a,
			Read:    false,
			Later:   false,
			Content: nil,
		})
	}
	feed = Feed{
		UserID:     userID,
		ID:         bson.NewObjectId().Hex(),
		PublicID:   publicID,
		URL:        url,
		Title:      title,
		Categories: []string{categoryID},
		Articles:   articles,
	}
	bytes, err := json.Marshal(&feed)
	if err != nil {
		return feed, err
	}
	_, err = stmt.Exec(string(bytes))
	return feed, err
}

func (m *FeedModel) UpdateArticles(userID, feedID string, articles []Article) (err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(userID, feedID)
	if err != nil {
		return err
	}
	feed := Feed{}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return errors.New("not_found")
	}
	for k := range articles {
		if articles[k].Later {
			for _, oldV := range feed.Articles {
				if oldV.URL == articles[k].URL {
					articles[k].Content = oldV.Content
					break
				}
			}
		}
	}
	stmt, err = m.DB.Prepare(`UPDATE feed SET json = JSON_SET(json, "$.articles", ?) WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(&articles)
	_, err = stmt.Exec(string(bytes), userID, feedID)
	return nil
}

func (m *FeedModel) GetFeedByID(userID, feedID string) (feed Feed, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return feed, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return feed, err
	}
	rows, err := stmt.Query(userID, feedID)
	if err != nil {
		return feed, err
	}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return feed, err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return feed, errors.New("not_found")
	}
	return feed, err
}

func (m *FeedModel) GetFeedByPublicID(userID, publicID string) (feed Feed, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(publicID) {
		return feed, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.pid") = ?`)
	if err != nil {
		return feed, err
	}
	rows, err := stmt.Query(userID, publicID)
	if err != nil {
		return feed, err
	}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return feed, err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return feed, errors.New("not_found")
	}
	return feed, err
}

func (m *FeedModel) GetFeedsByCategoryID(userID, categoryID string) (feeds []Feed, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(categoryID) {
		return feeds, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ?`)
	if err != nil {
		return feeds, err
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return feeds, err
	}
	var allFeeds []Feed
	feed := Feed{} 
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return feeds, err
		}
		err = json.Unmarshal(bytes, &feed)
		allFeeds = append(allFeeds, feed)
	} else {
		return feeds, errors.New("not_found")
	}
	for _, feed = range allFeeds {
		categories := feed.Categories
		for _, cid := range categories {
			if cid == categoryID {
				feeds = append(feeds, feed)
				break
			}
		}
	}
	return feeds, err
}

func (m *FeedModel) EditFeed(userID, feedID, title string, categoryIDs []string) (feed Feed, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return feed, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`UPDATE feed SET json = JSON_SET(json, "$.title", ?, "$.categories", ?) WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return feed, err
	}
	bytes, err := json.Marshal(&categoryIDs)
	res, err := stmt.Exec(title, string(bytes), userID, feedID)
	count, _ := res.RowsAffected()
	if count == 0 {
		return feed, errors.New("not_found")
	}
	stmt, err = m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return feed, err
	}
	rows, err := stmt.Query(userID, feedID)
	if err != nil {
		return feed, err
	}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return feed, err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return feed, errors.New("not_found")
	}
	return feed, err
}

func (m *FeedModel) RemoveFeed(userID, feedID string) (err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`DELETE FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(userID, feedID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return errors.New("not_found")
	}
	return nil
}

func (m *FeedModel) GetArticleByURL(userID, feedID, url string) (article Article, err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return article, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return article, err
	}
	rows, err := stmt.Query(userID, feedID)
	if err != nil {
		return article, err
	}
	feed := Feed{}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return article, err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return article, errors.New("not_found")
	}
	for k := range feed.Articles {
		if feed.Articles[k].URL == url {
			article = feed.Articles[k]
			break
		}
	}
	return article, err
}

func (m *FeedModel) GetLaterArticle(userID string) (articles []Article, err error) {
	if !bson.IsObjectIdHex(userID){
		return articles, errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ?`)
	if err != nil {
		return articles, err
	}
	rows, err := stmt.Query(userID)
	if err != nil {
		return articles, err
	}
	var allFeeds []Feed
	feed := Feed{} 
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return articles, err
		}
		err = json.Unmarshal(bytes, &feed)
		allFeeds = append(allFeeds, feed)
	} else {
		return articles, errors.New("not_found")
	}
	for _, feed = range allFeeds {
		allArticles := feed.Articles
		for _, article := range allArticles {
			if article.Later {
				articles = append(articles, article)
			}
		}
	}
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Content.Published >= articles[j].Content.Published
	})
	return articles, nil
}

func (m *FeedModel) EditArticle(userID, feedID, url string, read, later bool, article PublicArticle) (err error) {
	if !bson.IsObjectIdHex(userID) || !bson.IsObjectIdHex(feedID) {
		return errors.New("not_id")
	}
	stmt, err := m.DB.Prepare(`SELECT * FROM feed WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return err
	}
	rows, err := stmt.Query(userID, feedID)
	if err != nil {
		return err
	}
	feed := Feed{}
	if rows.Next() {
		var bytes []byte
		err = rows.Scan(&bytes)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bytes, &feed)
	} else {
		return errors.New("not_found")
	}
	for k := range feed.Articles {
		if feed.Articles[k].URL == url {
			feed.Articles[k].Read = read
			feed.Articles[k].Later = later
			if later == true {
				feed.Articles[k].Content = &article
			} else {
				feed.Articles[k].Content = nil
			}
			break
		}
	}
	stmt, err = m.DB.Prepare(`UPDATE feed SET json = JSON_SET(json, "$.articles", ?) WHERE JSON_EXTRACT(json, "$.uid") = ? AND JSON_EXTRACT(json, "$.id") = ?`)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(&feed.Articles)
	_, err = stmt.Exec(string(bytes), userID, feedID)
	return nil
}


