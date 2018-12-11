package app

import (
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/resolver"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/go-chi/chi"
	"net/http"
	"os"
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
	clientKey := os.Getenv("ClientKey")
	if clientKey == "" {
		fmt.Println("Client Key Not Found")
		return nil
	}
	return App(Config{
		Violet: violetSdk.Config{
			ClientID:   "5c0d08b28cb2530707c27f50",
			ClientKey:  clientKey,
			ServerHost: "https://oauth.xmatrix.studio/api/v2/",
			LoginURL:   "https://oauth.xmatrix.studio/Verify/Authorize/",
		},
	})
}
