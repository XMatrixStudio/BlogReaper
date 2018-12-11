package resolver

import (
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/service"
	"github.com/kataras/go-sessions"
)

type Resolver struct {
	Service *service.Service
	Session *sessions.Session
}

var resolver *Resolver

func DefaultResolver() *Resolver {
	if resolver == nil {
		resolver = &Resolver{
			Service: service.NewService(),
			Session: nil,
		}
	}
	return resolver
}

type mutationResolver struct{ *Resolver }

func (r *Resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}
