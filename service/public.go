package service

import (
	"encoding/xml"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/kataras/iris/core/errors"
	"io/ioutil"
	"net/http"
)

type PublicService interface {
	GetPublicFeed(url string) (publicFeed model.PublicFeed, err error)
	UpdatePublicFeed(url string) (publicFeed model.PublicFeed, err error)
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

// 从数据库中获取PublicFeed
func (s *publicService) GetPublicFeed(url string) (publicFeed model.PublicFeed, err error) {
	panic("not implement")
}

type feed struct {
	Title    string  `xml:"title"`
	Subtitle string  `xml:"subtitle"`
	Author   author  `xml:"author"`
	Entrys   []entry `xml:"entry"`
}

type author struct {
	Name string `xml:"name"`
}

type entry struct {
	Title     string `xml:"title"`
	Link      link   `xml:"link"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Content   string `xml:"content"`
	Summary   string `xml:"summary"`
}

type link struct {
	Href string `xml:"href,attr"`
}

type category struct {
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
	xml.Unmarshal(bytes, &publicFeed)
	panic("not implement")
	//// TODO
	//// 获取atom.xml
	//client := http.DefaultClient
	//resp, err := client.Get(url)
	//if err != nil {
	//	return nil, errors.New("http_request_fail")
	//}
	//defer resp.Body.Close()
	//con, _ := ioutil.ReadAll(resp.Body)
	//var result feed
	//// 解析atom.xml
	//err = xml.Unmarshal(con, &result)
	//fmt.Println(result.Entrys)
}
