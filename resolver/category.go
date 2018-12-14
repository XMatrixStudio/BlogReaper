package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *categoryResolver) Feeds(ctx context.Context, obj *graphql.Category, id *string) ([]graphql.Feed, error) {
	if obj == nil {
		return nil, nil
	}
	if id == nil {
		return obj.Feeds, nil
	}
	for _, v := range obj.Feeds {
		if v.ID == *id {
			return []graphql.Feed{v}, nil
		}
	}
	return nil, errors.New("invalid_id")
}
