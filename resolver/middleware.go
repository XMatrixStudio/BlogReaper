package resolver

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/kataras/go-sessions"
	"net/http"
	"time"
)

var sessionsManager = sessions.New(sessions.Config{
	Cookie:                      "BlogReaperSession",
	Expires:                     time.Hour * 24,
	DisableSubdomainPersistence: false,
})

func SessionHttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "session", sessionsManager.Start(w, r))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SessionResolverMiddleware(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	resolver.Session = ctx.Value("session").(*sessions.Session)
	return next(ctx)
}
