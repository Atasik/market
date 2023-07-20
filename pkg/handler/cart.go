package handler

import (
	"encoding/json"
	"io/ioutil"
	"market/pkg/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type cartInput struct {
	Amount int `json:"amount" validate:"required"`
}

func (h *Handler) addProductToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var input cartInput

	err = json.Unmarshal(body, &input)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(input)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	Cart, err := h.Services.Cart.GetByUserID(sess.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Cart.AddProduct(Cart.ID, sess.ID, productID, input.Amount)
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

func (h *Handler) deleteProductFromCart(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getProductsFromCart(w http.ResponseWriter, r *http.Request) {
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
