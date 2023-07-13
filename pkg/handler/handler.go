package handler

import (
	"html/template"
	"market/pkg/service"
	"market/pkg/session"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	Sessions *session.SessionsManager
	Services *service.Service
}

func (h *Handler) InitRoutes() *mux.Router {
	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static")),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", h.Index).Methods("GET")
	r.HandleFunc("/about", h.About).Methods("GET")
	r.HandleFunc("/history", h.History).Methods("GET")

	r.HandleFunc("/products/new", h.CreateProductForm).Methods("GET")
	r.HandleFunc("/products/new", h.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}/reviews/new", h.AddReview).Methods("POST")
	r.HandleFunc("/products/{id}/reviews/update", h.UpdateReview).Methods("POST")
	r.HandleFunc("/products/{id}/reviews/delete", h.DeleteReview).Methods("DELETE")
	r.HandleFunc("/products/update/{id}", h.UpdateProductForm).Methods("GET")
	r.HandleFunc("/products/update/{id}", h.UpdateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", h.GetProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", h.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", h.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/basket/{id}", h.AddProductToBasket).Methods("GET")
	r.HandleFunc("/basket/{id}", h.DeleteProductFromBasket).Methods("DELETE")
	r.HandleFunc("/basket", h.Basket).Methods("GET")
	r.HandleFunc("/register_order", h.CreateOrder).Methods("GET")

	r.HandleFunc("/register", h.Register).Methods("GET")
	r.HandleFunc("/login", h.Login).Methods("GET")
	r.HandleFunc("/logout", h.Logout).Methods("GET")

	r.HandleFunc("/sign_up", h.SignUp).Methods("POST")
	r.HandleFunc("/sign_in", h.SignIn).Methods("POST")

	r.PathPrefix("/static/").Handler(staticHandler)

	return r
}
