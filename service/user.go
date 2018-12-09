package service

import (
	"github.com/XMatrixStudio/BlogReaper/model"
	"github.com/XMatrixStudio/Violet.SDK.Go"
)

type UserService interface {
	InitViolet(c violetSdk.Config)
	GetLoginURL(backURL string) (url, state string)
	LoginByCode(code string) (userID string, err error)
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
	// 获取用户Token
	res, err := s.Violet.GetToken(code)
	if err != nil {
		return
	}
	// 保存数据并获取用户信息
	user, err := s.Model.GetUserByID(res.UserID)
	userID = user.VioletID.Hex()
	if err == nil { // 数据库已存在该用户
		s.Model.SetUserToken(user.VioletID.Hex(), res.Token)
	} else if err.Error() == "not_found" { // 数据库不存在此用户
		userNew, err := s.Violet.GetUserBaseInfo(res.UserID, res.Token)
		if err != nil {
			return userID, err
		}
		err = s.Model.AddUser(res.UserID, res.Token, userNew.Email, userNew.Name, userNew.Info.Avatar, userNew.Info.Bio, userNew.Info.Gender)
	}
	return
}
