package handler

import (
	"encoding/json"
	"io"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
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
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var review model.Review

	err = json.Unmarshal(body, &review)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(review)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	review.ProductID = productID
	review.UserID = sess.UserID
	review.Username = sess.Username

	review.ID, err = h.Services.Review.Create(review)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Review was created with id LastInsertId: %v", review.ID)

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews, err = h.Services.Review.GetAll(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		Limit:     5,
		Offset:    0,
		ProductID: productID,
		SortBy:    model.SortByViews,
		SortOrder: model.DESCENDING,
	}

	product.RelatedProducts, err = h.Services.Product.GetProductsByCategory(product.Category, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateReview(w http.ResponseWriter, r *http.Request) {
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
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var input model.UpdateReviewInput

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

	currentTime := time.Now()
	input.UpdatedAt = &currentTime

	err = h.Services.Review.Update(sess.UserID, productID, input)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Review by userID [%v] to productID [%v] was updated: %v", sess.UserID, productID)

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews, err = h.Services.Review.GetAll(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		Limit:     5,
		Offset:    0,
		ProductID: productID,
		SortBy:    model.SortByViews,
		SortOrder: model.DESCENDING,
	}

	product.RelatedProducts, err = h.Services.Product.GetProductsByCategory(product.Category, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["reviewId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	err = h.Services.Review.Delete(sess.UserID, reviewID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	h.Logger.Infof("Review by userID [%v] to productID [%v] was deleted", sess.UserID, productID)

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews, err = h.Services.Review.GetAll(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		Limit:     5,
		Offset:    0,
		ProductID: productID,
		SortBy:    model.SortByViews,
		SortOrder: model.DESCENDING,
	}

	product.RelatedProducts, err = h.Services.Product.GetProductsByCategory(product.Category, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}
