package service

import "github.com/XMatrixStudio/BlogReaper/model"

type Service struct {
	User UserService
}

func NewService() *Service {
	s, m := &Service{}, model.NewModel()
	s.User = NewUserService(s, &model.UserModel{m})
	return s
}
