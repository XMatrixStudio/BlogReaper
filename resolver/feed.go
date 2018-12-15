package resolver

import (
	"context"
	"errors"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *feedResolver) Articles(ctx context.Context, obj *graphql.Feed, page *int, numPerPage *int) ([]graphql.Article, error) {
	if obj == nil {
		return nil, nil
	}
	if page == nil && numPerPage == nil {
		return obj.Articles, nil
	}
	if page == nil || numPerPage == nil || *page <= 0 || *numPerPage <= 0 {
		return nil, errors.New("invalid_params")
	}
	start := (*page - 1) * (*numPerPage)
	end := (*page-1)*(*numPerPage) + *numPerPage
	if len(obj.Articles) < start {
		return nil, nil
	} else if len(obj.Articles) > end {
		return obj.Articles[start:end], nil
	} else {
		return obj.Articles[start:], nil
	}
}
