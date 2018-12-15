package service

import (
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/model"
)

type CategoryService interface {
	GetModel() *model.CategoryModel
	AddCategory(userID, name string) (category graphql.Category, err error)
	GetCategories(userID string) (categories []graphql.Category, err error)
	EditCategory(userID, categoryID, newName string) (success bool, err error)
	RemoveCategory(userID, categoryID string) (success bool, err error)
}

type categoryService struct {
	Model   *model.CategoryModel
	Service *Service
}

func NewCategoryService(s *Service, m *model.CategoryModel) CategoryService {
	return &categoryService{
		Model:   m,
		Service: s,
	}
}

func (s *categoryService) GetModel() *model.CategoryModel {
	return s.Model
}

func (s *categoryService) AddCategory(userID, name string) (category graphql.Category, err error) {
	c, err := s.Model.AddCategory(userID, name)
	if err != nil {
		return category, err
	}
	category = graphql.Category{
		ID:    c.ID.Hex(),
		Name:  c.Name,
		Feeds: nil,
	}
	return
}

func (s *categoryService) GetCategories(userID string) (categories []graphql.Category, err error) {
	cs, err := s.Model.GetCategories(userID)
	if err != nil {
		return categories, err
	}
	for _, c := range cs {
		feeds, err := s.Service.Feed.GetFeedsByCategoryID(userID, c.ID.Hex())
		if err != nil && err.Error() == "not_found" {
			feeds = nil
		}
		categories = append(categories, graphql.Category{
			ID:    c.ID.Hex(),
			Name:  c.Name,
			Feeds: feeds,
		})
	}
	return
}

func (s *categoryService) EditCategory(userID, categoryID, newName string) (success bool, err error) {
	category, err := s.Model.GetCategoryByName(userID, newName)
	if err != nil {
		return false, err
	}
	if category.Name != "" && category.ID.Hex() != categoryID {
		return false, errors.New("repeat_name")
	}
	if category.Name != "" {
		return true, nil
	}
	return s.Model.EditCategory(userID, categoryID, newName)
}

func (s *categoryService) RemoveCategory(userID, categoryID string) (success bool, err error) {
	panic("not implement")
}
