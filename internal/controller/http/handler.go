package http

import (
	v1 "market/internal/controller/http/v1"
	"market/internal/service"
	"market/pkg/auth"
	"net/http"

	_ "market/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	logger       *zap.SugaredLogger
	services     *service.Service
	tokenManager auth.TokenManager
	validator    *validator.Validate
}

func NewHandler(services *service.Service, validator *validator.Validate, logger *zap.SugaredLogger, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		validator:    validator,
		logger:       logger,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	h.initAPI(r)

	return r
}

func (h *Handler) initAPI(router *mux.Router) {
	handlerV1 := v1.NewHandler(h.services, h.validator, h.logger, h.tokenManager)
	api := router.PathPrefix("/api").Subrouter()
	handlerV1.Init(api)
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	h.initAPI(r)

	mux := h.accessLogMiddleware(r)
	mux = panicMiddleware(mux)

	return mux
}

//endpoints
// /api/v1/products - GET
// /api/v1/products/category/{categoryName} - GET
// /api/v1/product/{productId} - GET
// /api/v1/product/{productId} - DELETE
// /api/v1/product/{productId} - POST
// /api/v1/product/{productId} - PUT
// /api/v1/product/{productId}/review - POST
// /api/v1/product/{productId}/review/{reviewId} - PUT
// /api/v1/product/{productId}/review/{reviewId} - POST

// /api/v1/cart - GET
// /api/v1/cart - DELETE
// /api/v1/cart/{productId} - PUT
// /api/v1/cart/{productId} - POST
// /api/v1/cart/{productId} - DELETE

// /api/v1/orders - GET
// /api/v1/order - GET
// /api/v1/order - POST

// /api/v1/user/sign-up - POST
// /api/v1/user/sign-in - POST
// /api/v1/user/{userId}/products - GET
