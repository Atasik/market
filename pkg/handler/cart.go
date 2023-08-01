package handler

import (
	"encoding/json"
	"io"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// @Summary Add product to cart
// @Security ApiKeyAuth
// @Tags cart
// @ID add-product-to-cart
// @Accept json
// @Product json
// @Param   productId path integer true "ID of product to add to cart"
// @Param input body cartInput true "Amount of products"
// @Success 200 {object} getProductsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/cart/{productId} [post]
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

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			Offset:    defaultOffset,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary Get products from cart
// @Security ApiKeyAuth
// @Tags cart
// @ID get-products-from-cart
// @Product json
// @Param   sort_by query   string false "sort by" Enums(created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success 200 {object} getProductsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/cart [get]
func (h *Handler) getProductsFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart, err := h.Services.Cart.GetByUserID(sess.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID, q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

type cartInput struct {
	Amount int `json:"amount" validate:"required"`
}

// @Summary Update product amount from cart
// @Security ApiKeyAuth
// @Tags cart
// @ID update-product-amount-from-cart
// @Accept json
// @Product json
// @Param   productId path integer true "ID of product to update"
// @Param input body cartInput true "Amount of products"
// @Success 200 {object} getProductsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/cart/{productId} [put]
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

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     0,
			Offset:    0,
			SortBy:    model.SortByDate,
			SortOrder: model.DESCENDING,
		},
	}

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary Delete product from cart
// @Security ApiKeyAuth
// @Tags cart
// @ID delete-product-from-cart
// @Product json
// @Param   productId path integer true "ID of product to delete"
// @Success 200 {object} getProductsResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/cart/{productId} [delete]
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

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     0,
			Offset:    0,
			SortBy:    model.SortByDate,
			SortOrder: model.DESCENDING,
		},
	}

	products, err := h.Services.Cart.GetAllProducts(sess.UserID, cart.ID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary Delete products from cart
// @Security ApiKeyAuth
// @Tags cart
// @ID delete-products-from-cart
// @Product json
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/cart [delete]
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
