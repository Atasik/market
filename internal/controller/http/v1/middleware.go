package v1

import (
	"context"
	"errors"
	"fmt"
	"market/internal/model"
	"market/pkg/auth"
	"net/http"
	"strconv"
	"strings"
)

type OptionsKey string

const (
	OptionsContextKey OptionsKey = "query_options"
)

const (
	authorizationHeader = "Authorization"
	defaultSortField    = "created_at"
	defaultPage         = 1
	defaultLimit        = 25
	maxLimit            = 50
)

var (
	ErrNoQuery = errors.New("no query")
)

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middleware", r.URL.Path)

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

		token, err := h.tokenManager.Parse(headerParts[1])
		if err != nil {
			newErrorResponse(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := auth.ContextWithToken(r.Context(), token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func queryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("queryMiddleware", r.URL.Path)

		sortBy := strings.ToLower(r.URL.Query().Get("sort_by"))
		sortOrder := strings.ToUpper(r.URL.Query().Get("sort_order"))
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			limit = defaultLimit
		}
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = defaultPage
		}

		if sortBy == "" {
			sortBy = defaultSortField
		}

		if sortOrder == "" {
			sortOrder = model.DESCENDING
		}

		if limit < 1 {
			limit = defaultLimit
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		if page < 1 {
			page = defaultPage
		}

		options := &Options{
			SortBy:    sortBy,
			SortOrder: sortOrder,
			Limit:     limit,
			Offset:    (page - 1) * limit,
		}
		ctx := contextWithOptions(r.Context(), options)

		next(w, r.WithContext(ctx))
	}
}

type Options struct {
	SortBy, SortOrder string
	Limit             int
	Offset            int
}

func optionsFromContext(ctx context.Context) (*Options, error) {
	options, ok := ctx.Value(OptionsContextKey).(*Options)
	if !ok || options == nil {
		return nil, ErrNoQuery
	}
	return options, nil
}

func contextWithOptions(ctx context.Context, opts *Options) context.Context {
	return context.WithValue(ctx, OptionsContextKey, opts)
}
