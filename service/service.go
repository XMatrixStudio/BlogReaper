package service

import (
	"fmt"
	"github.com/XMatrixStudio/BlogReaper/model"
	"os"
)

type Service struct {
	User   UserService
	Feed   FeedService
	Public PublicService
}

func NewService() *Service {
	s := &Service{}
	m, err := model.NewModel()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
	s.User = NewUserService(s, &model.UserModel{Model: &model.Model{BucketName: "user", DB: m.DB}})
	s.Feed = NewFeedService(s, &model.FeedModel{Model: &model.Model{BucketName: "feed", DB: m.DB}})
	s.Public = NewPublicService(s, &model.PublicModel{Model: &model.Model{BucketName: "public", DB: m.DB}})
	return s
}
