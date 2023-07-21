package handler

import (
	"encoding/json"
	"io"
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

	body, err := io.ReadAll(r.Body)
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

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Cart.AddProduct(cart.ID, sess.UserID, productID, input.Amount)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Product was added to cart with id LastInsertId: %v", productID)

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
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

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) updateProductAmountFromCart(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(r.Body)
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

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Services.Cart.UpdateProductAmount(cart.ID, sess.UserID, productID, input.Amount)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Product from Cart %v was updated: %v %v", cart.ID, productID)

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
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
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Services.Cart.DeleteProduct(cart.ID, sess.UserID, productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Product from Cart %v was deleted: %v %v", cart.ID, productID)

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) deleteProductsFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Services.Cart.DeleteAllProducts(cart.ID, sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Products was deleted from Cart: %v %v", cart.ID)

	newStatusReponse(w, "done", http.StatusOK)
}
