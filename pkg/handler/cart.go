package handler

import (
	"market/pkg/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	Cart, err := h.Services.Cart.GetByUserID(sess.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Cart.AddProduct(Cart.ID, sess.ID, productId)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products, err := h.Services.Cart.GetProducts(sess.ID, Cart.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) DeleteProductFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	Cart, err := h.Services.Cart.GetByUserID(sess.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Cart.DeleteProduct(Cart.ID, sess.ID, productId)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products, err := h.Services.Cart.GetProducts(sess.ID, Cart.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	Cart, err := h.Services.Cart.GetByUserID(sess.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products, err := h.Services.Cart.GetProducts(sess.ID, Cart.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}
