package main

import (
	"flag"
	"github.com/99designs/gqlgen/handler"
	"github.com/XMatrixStudio/BlogReaper/graphql"
	"github.com/XMatrixStudio/BlogReaper/resolver"
	"github.com/go-chi/chi"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 加载配置文件
	configFile := flag.String("c", "config/config.yaml", "Where is your config file?")
	flag.Parse()
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Can't find the config file in %v", *configFile)
		return
	}
	log.Printf("Load the config file in %v", *configFile)
	conf := resolver.Config{}
	yaml.Unmarshal(data, &conf)
	log.Fatal(http.ListenAndServe(":30038", app(conf)))
}

func app(c resolver.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(resolver.SessionHttpMiddleware)
	r.Handle("/", handler.Playground("BlogReaper", "/graphql"))
	r.Handle("/graphql", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{
		Resolvers: resolver.DefaultResolver(c),
	}), handler.ResolverMiddleware(resolver.SessionResolverMiddleware)))
	return r
}
