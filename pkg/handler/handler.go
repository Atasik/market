package handler

import (
	"market/pkg/service"
	"net/http"

	_ "market/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const (
	appJSON       = "application/json"
	multiFormData = "multipart/form-data"
)

type Handler struct {
	Logger    *zap.SugaredLogger
	Services  *service.Service
	Validator *validator.Validate
}

func (h *Handler) InitRoutes() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/api/products", query(h.getAllProducts)).Methods("GET")
	r.HandleFunc("/api/category/{categoryName}", query(h.getProductsByCategory)).Methods("GET")
	r.HandleFunc("/api/products/{userId}", query(h.getProductsByUserID)).Methods("GET")

	authR := mux.NewRouter()

	authR.HandleFunc("/api/orders", h.getOrders).Methods("GET")
	authR.HandleFunc("/api/order", h.createOrder).Methods("GET")
	authR.HandleFunc("/api/order/{orderId}", h.getOrder).Methods("GET")

	authR.HandleFunc("/api/product", h.createProduct).Methods("POST")
	authR.HandleFunc("/api/product/{productId}", h.updateProduct).Methods("PUT")
	r.HandleFunc("/api/product/{productId}", h.getProductByID).Methods("GET")
	authR.HandleFunc("/api/product/{productId}", h.deleteProduct).Methods("DELETE")
	authR.HandleFunc("/api/product/{productId}/addReview", h.createReview).Methods("POST")
	authR.HandleFunc("/api/product/{productId}/updateReview/{reviewId}", h.updateReview).Methods("PUT")
	authR.HandleFunc("/api/product/{productId}/deleteReview/{reviewId}", h.deleteReview).Methods("DELETE")

	authR.HandleFunc("/api/cart", h.getProductsFromCart).Methods("GET")
	authR.HandleFunc("/api/cart/{productId}", h.updateProductAmountFromCart).Methods("PUT")
	authR.HandleFunc("/api/cart/{productId}", h.addProductToCart).Methods("POST")
	authR.HandleFunc("/api/cart/{productId}", h.deleteProductFromCart).Methods("DELETE")
	authR.HandleFunc("/api/cart", h.deleteProductsFromCart).Methods("DELETE")

	r.HandleFunc("/api/register", h.signUp).Methods("POST")
	r.HandleFunc("/api/login", h.signIn).Methods("POST")

	r.PathPrefix("/api/").Handler(auth(h.Services.User, authR))

	mux := accessLog(h.Logger, r)
	mux = panic(mux)

	return mux
}
