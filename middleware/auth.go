package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"market/pkg/session"
)

var (
	noAuthUrls = map[string]struct{}{
		"/login":    {},
		"/register": {},
		"/sign_in":  {},
		"/sign_up":  {},
	}
	noSessUrls = map[string]struct{}{
		"/":        {},
		"/product": {},
		"/about":   {},
		"/privacy": {},
		"/basket":  {},
	}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)

		_, canBeWithouthSess := noSessUrls[r.URL.Path]
		if err != nil && !canBeWithouthSess && !strings.HasPrefix(r.URL.Path, "/static") && !strings.HasPrefix(r.URL.Path, "/products") {
			fmt.Println("no auth")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
