package handler

import (
	"fmt"
	"market/pkg/service"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	authorizationHeader = "Authorization"
)

func accessLog(logger *zap.SugaredLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("access log middleware")
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infow("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"time", time.Since(start),
		)
	})
}

func auth(s service.User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware")

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			newErrorResponse(w, "empty auth header", http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			newErrorResponse(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		if len(headerParts[1]) == 0 {
			newErrorResponse(w, "empty token", http.StatusUnauthorized)
			return
		}

		sess, err := s.CheckToken(headerParts[1])
		if err != nil {
			newErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := service.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func panic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
