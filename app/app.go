package app

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/resolver"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/go-chi/chi"
	"net/http"
)

func App(c Config) http.Handler {
	r := chi.NewRouter()
	r.Use(resolver.SessionHttpMiddleware)
	dr := resolver.DefaultResolver()
	dr.Service.User.InitViolet(c.Violet)
	r.Handle("/", handler.Playground("BlogReaper", "/graphql"))
	r.Handle("/graphql", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: dr,
	}), handler.ResolverMiddleware(resolver.SessionResolverMiddleware)))
	return r
}

func TestApp() http.Handler {
	return App(Config{
		Violet: violetSdk.Config{},
	})
}
