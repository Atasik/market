package handler

import (
	"market/pkg/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	appJSON       = "application/json"
	multiFormData = "multipart/form-data"
)

type Handler struct {
	Logger   *zap.SugaredLogger
	Services *service.Service
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.GetProducts).Methods("GET")
	r.HandleFunc("/api/about", h.About).Methods("GET")

	authR := mux.NewRouter()

	authR.HandleFunc("/api/orders", h.GetOrders).Methods("GET")
	authR.HandleFunc("/api/product", h.CreateProduct).Methods("POST")
	authR.HandleFunc("/api/product/{productId}", h.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/product/{productId}", h.GetProduct).Methods("GET")
	authR.HandleFunc("/api/product/{productId}", h.DeleteProduct).Methods("DELETE")
	authR.HandleFunc("/api/product/{productId}/addReview", h.CreateReview).Methods("POST")
	authR.HandleFunc("/api/product/{productId}/updateReview/{reviewId}", h.UpdateReview).Methods("PUT")
	authR.HandleFunc("/api/product/{productId}/deleteReview/{reviewId}", h.DeleteReview).Methods("DELETE")

	authR.HandleFunc("/api/cart", h.GetProductsFromCart).Methods("GET")
	authR.HandleFunc("/api/cart/{productId}", h.AddProductToCart).Methods("GET")
	authR.HandleFunc("/api/cart/{productId}", h.DeleteProductFromCart).Methods("DELETE")

	authR.HandleFunc("/api/createOrder", h.CreateOrder).Methods("GET")

	r.HandleFunc("/api/register", h.SignUp).Methods("POST")
	r.HandleFunc("/api/login", h.SignIn).Methods("POST")

	r.PathPrefix("/api/").Handler(Auth(h.Services.User, authR))

	return r
}
