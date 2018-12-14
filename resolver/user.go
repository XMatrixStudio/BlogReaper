package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *userResolver) Categories(ctx context.Context, obj *graphql.User, id *string) ([]graphql.Category, error) {
	if obj == nil {
		return nil, nil
	}
	if id == nil {
		return obj.Categories, nil
	}
	for _, v := range obj.Categories {
		if v.ID == *id {
			return []graphql.Category{v}, nil
		}
	}
	return nil, errors.New("invalid_id")
}
