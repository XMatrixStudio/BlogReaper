package resolver

import (
	"context"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/service"
	"github.com/kataras/go-sessions"
)

type Resolver struct {
	User    service.UserService
	Session *sessions.Session
}

var resolver *Resolver

func DefaultResolver(c Config) *Resolver {
	if resolver == nil {
		resolver = &Resolver{
			User:    service.NewUserService(),
			Session: nil,
		}
		resolver.User.InitViolet(c.Violet)
	}
	return resolver
}

func (r *Resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context) (*graphql.User, error) {
	panic("not implemented")
}
