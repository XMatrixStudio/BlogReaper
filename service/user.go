package service

import "github.com/XMatrixStudio/Violet.SDK.Go"

type UserService interface {
	InitViolet(c violetSdk.Config)
	GetLoginURL(backURL string) (url, state string)
	LoginByCode(code string) (userID string, err error)
}

type userService struct {
	Violet violetSdk.Violet
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) InitViolet(c violetSdk.Config) {
	s.Violet = violetSdk.NewViolet(c)
}

func (s *userService) GetLoginURL(backUrl string) (url, state string) {
	return s.Violet.GetLoginURL(backUrl)
}

func (s *userService) LoginByCode(code string) (userID string, err error) {
	// 获取用户Token
	_, err = s.Violet.GetToken(code)
	if err != nil {
		return
	}
	// TODO
	return
}
