package service

import "github.com/XMatrixStudio/BlogReaper/model"

type PublicService interface {
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
