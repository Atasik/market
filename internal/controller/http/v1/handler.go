package v1

import (
	"market/internal/service"
	"market/pkg/auth"
	"market/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	appJSON       = "application/json"
	multiFormData = "multipart/form-data"
)

type Handler struct {
	logger       logger.Logger
	services     *service.Service
	tokenManager auth.TokenManager
	validator    *validator.Validate
}

func NewHandler(services *service.Service, validator *validator.Validate, logger logger.Logger, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		validator:    validator,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *mux.Router) {
	r := api.PathPrefix("/v1").Subrouter()
	h.initCartRoutes(r)
	h.initProductRoutes(r)
	h.initProductsRoutes(r)
	h.initOrderRoutes(r)
	h.initOrdersRoutes(r)
	h.initUserRoutes(r)
}
