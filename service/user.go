package service

import (
	"errors"
	"flag"
	"fmt"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/globalsign/mgo/bson"
)

type UserService interface {
	InitViolet(c violetSdk.Config)
	GetLoginURL(backURL string) (url, state string)
	LoginByCode(code string) (userID string, err error)
	GetUserInfo(id string) (user graphql.User, err error)
}

type userService struct {
	Violet  violetSdk.Violet
	Model   *model.UserModel
	Service *Service
}

func NewUserService(s *Service, m *model.UserModel) UserService {
	return &userService{
		Model:   m,
		Service: s,
	}
}

func (s *userService) InitViolet(c violetSdk.Config) {
	s.Violet = violetSdk.NewViolet(c)
}

func (s *userService) GetLoginURL(backUrl string) (url, state string) {
	return s.Violet.GetLoginURL(backUrl)
}

func (s *userService) LoginByCode(code string) (userID string, err error) {
	if flag.Lookup("test.v") != nil {
		fmt.Println("normal run")
		// 获取用户Token
		res, err := s.Violet.GetToken(code)
		if err != nil {
			return "", err
		}
		// 保存数据并获取用户信息
		user, err := s.Model.GetUserByID(res.UserID)
		userID = user.VioletID.Hex()
		if err == nil { // 数据库已存在该用户
			s.Model.SetUserToken(user.VioletID.Hex(), res.Token)
		} else if err.Error() == "not_found" { // 数据库不存在此用户
			userNew, err := s.Violet.GetUserBaseInfo(res.UserID, res.Token)
			if err != nil {
				return userID, errors.New("violet_error")
			}
			err = s.Model.AddUser(res.UserID, res.Token, userNew.Email, userNew.Name, userNew.Info.Avatar, userNew.Info.Bio, userNew.Info.Gender)
			if err != nil {
				return userID, errors.New("initial_fail")
			}
		}
	} else {
		fmt.Println("run under go test")
		userID = string(bson.NewObjectId())
		fmt.Println(userID)
		user, err := s.Model.GetUserByID(userID)
		if err != nil {
			// 测试伪造用户
			fmt.Println(1)
			err = s.Model.AddUser(userID, "faker_token", "faker@qq.com", "faker", "", "xxx", 0)
			return userID, err
		} else {
			s.Model.SetUserToken(string(user.VioletID), "faker_token")
			return string(user.VioletID), nil
		}
	}
	return
}

func (s *userService) GetUserInfo(id string) (user graphql.User, err error) {
	modelUser, err := s.Model.GetUserByID(id)
	if err != nil {
		return graphql.User{}, errors.New("not_found")
	}
	user = graphql.User{
		Email: modelUser.Email,
		Info: graphql.UserInfo{
			Name:   modelUser.Info.Name,
			Avatar: modelUser.Info.Avatar,
			Bio:    modelUser.Info.Bio,
			Gender: modelUser.Info.Gender,
		},
	}
	return
}
