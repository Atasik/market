package http

import (
	v1 "market/internal/controller/http/v1"
	"market/internal/service"
	"market/pkg/auth"
	"market/pkg/logger"
	"net/http"

	_ "market/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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

// endpoints
//
// /api/v1/products - GET - Получить продукты
//
// /api/v1/products/category/{categoryName} - GET - Получить продукты в категории
//
// /api/v1/product/{productId} - GET - Получить продукт
//
// /api/v1/product/{productId} - DELETE - Удалить продукт
//
// /api/v1/product/{productId} - POST - Создать продукт
//
// /api/v1/product/{productId} - PUT - Изменить продукт
//
// /api/v1/product/{productId}/review - POST - Создать отзыв
//
// /api/v1/product/{productId}/review/{reviewId} - PUT - Изменить отзыв
//
// /api/v1/product/{productId}/review/{reviewId} - DELETE - Создаить отзыв
//
// /api/v1/cart - GET - Получить содержимое корзины
//
// /api/v1/cart - DELETE - Очистить корзину
//
// /api/v1/cart/{productId} - PUT - Изменить кол-во продуктов в корзине
//
// /api/v1/cart/{productId} - POST - Добавить в корзину
//
// /api/v1/cart/{productId} - DELETE - Удалить из корзины
//
// /api/v1/orders - GET - Получить заказы
//
// /api/v1/order/{orderId} - GET - Получить заказ
//
// /api/v1/order - POST - Создать заказ
//
// /api/v1/user/sign-up - POST - Выйти
//
// /api/v1/user/sign-in - POST - Войти
//
// /api/v1/user/{userId}/products - GET - Получить продукты юзера
func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	h.initAPI(r)

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE", "PUT", "OPTIONS"})

	mux := h.accessLogMiddleware(r)
	mux = panicMiddleware(mux)
	return handlers.CORS(originsOk, headersOk, methodsOk)(mux)
}
