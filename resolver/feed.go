package resolver

import (
	"context"
	"github.com/XMatrixStudio/BlogReaper/graphql"
)

func (r *feedResolver) Articles(ctx context.Context, obj *graphql.Feed, page int, numPerPage int) ([]graphql.Article, error) {
	panic("not implemented")
}
