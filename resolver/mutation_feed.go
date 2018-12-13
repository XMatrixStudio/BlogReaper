package resolver

import (
	"context"
	"encoding/xml"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"io/ioutil"
	"net/http"
)

type feed struct {
	Title string  `xml:"title"`
	Subtitle string `xml:"subtitle"`
	Entrys []entry `xml:"entry"`
}

type entry struct {
	Title string     `xml:"title"`
	Link link 		 `xml:"link"`
	Href string      `xml:"href, attr"`
	Published string `xml:"published"`
	Updated string   `xml:"updated"`
	Content string   `xml:"content"`
	Summary string   `xml:"summary"`
}

type link struct {
	Href string `xml:"href,attr"`
}

func (r *mutationResolver) AddFeed(ctx context.Context, url string, categoryId *string, categoryName *string) (*graphql.Category, error) {
	if r.Session.Get("id") == nil {
		return nil, errors.New("not_login")
	}
	// TODO
	// 获取atom.xml
	client := http.DefaultClient
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.New("http_request_fail")
	}
	defer resp.Body.Close()
	con, _ := ioutil.ReadAll(resp.Body)
	var result feed
	// 解析atom.xml
	err = xml.Unmarshal(con, &result)

	return nil, nil
}

func (r *mutationResolver) EditArticle(ctx context.Context, url string, read *bool, later *bool) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}

	panic("not implemented")
}

func (r *mutationResolver) EditCategory(ctx context.Context, id string, name string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *mutationResolver) EditFeed(ctx context.Context, url string, title *string, categoryId *string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}

func (r *mutationResolver) RemoveFeed(ctx context.Context, url string) (bool, error) {
	if r.Session.Get("id") == nil {
		return false, errors.New("not_login")
	}
	panic("not implemented")
}
