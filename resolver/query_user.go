package resolver

import (
	"context"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/kataras/iris/core/errors"
)

func (r *queryResolver) User(ctx context.Context) (*graphql.User, error) {
	userID := r.Session.GetString("id")
	if userID == "" {
		return nil, errors.New("not_login")
	}
	user, err := r.Service.User.GetUserInfo(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
