package handler

import (
	"encoding/json"
	"io"
	"market/internal/model"
	"market/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type reviewInput struct {
	Text     string `json:"text" validate:"required"`
	Category string `json:"category" validate:"review_category,required"`
}

// @Summary	Create review
// @Security	ApiKeyAuth
// @Tags		review
// @ID			create-review
// @Accept		json
// @Product	json
// @Param		productId	path		integer			true	"ID of product for review"
// @Param		input		body		reviewInput	true	"Review content"
// @Success	201			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/product/{productId}/addReview [post]
func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
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

	var input reviewInput

	if err = json.Unmarshal(body, &input); err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	if err = h.validator.Struct(input); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var review model.Review
	review.Category = input.Category
	review.Text = input.Text
	review.ProductID = productID
	review.UserID = token.UserID
	review.Username = token.Username
	review.CreatedAt = time.Now()

	review.ID, err = h.services.Review.Create(review)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infof("Review was created with id LastInsertId: %v", review.ID)

	product, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviewQuery := model.ReviewQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	product.Reviews, err = h.services.Review.GetAll(productID, reviewQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productQuery := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     limitRelatedProducts,
			SortBy:    model.SortByViews,
			SortOrder: model.DESCENDING,
		},
		ProductID: productID,
	}

	product.RelatedProducts, err = h.services.Product.GetProductsByCategory(product.Category, productQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Update review
// @Security	ApiKeyAuth
// @Tags		review
// @ID			update-review
// @Accept		json
// @Product	json
// @Param		productId	path		integer			true	"ID of product"
// @Param		reviewId	path		integer			true	"ID of review"
// @Param		input		body		reviewInput	true	"Review content"
// @Success	201			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/product/{productId}/updateReview/{reviewId} [put]
func (h *Handler) updateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
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

	if err = json.Unmarshal(body, &input); err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	if err = h.validator.Struct(input); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	input.UpdatedAt = &currentTime

	if err = h.services.Review.Update(token.UserID, productID, input); err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infof("Review by userID [%v] to productID [%v] was updated: %v", token.UserID, productID)

	product, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviewQuery := model.ReviewQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	product.Reviews, err = h.services.Review.GetAll(productID, reviewQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productQuery := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     limitRelatedProducts,
			SortBy:    model.SortByViews,
			SortOrder: model.DESCENDING,
		},
		ProductID: productID,
	}

	product.RelatedProducts, err = h.services.Product.GetProductsByCategory(product.Category, productQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Delete review
// @Security	ApiKeyAuth
// @Tags		review
// @ID			delete-review
// @Product	json
// @Param		productId	path		integer	true	"ID of product"
// @Param		reviewId	path		integer	true	"ID of review"
// @Success	200			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/product/{productId}/deleteReview/{reviewId} [delete]
func (h *Handler) deleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["reviewId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	if err = h.services.Review.Delete(token.UserID, reviewID); err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
		return
	}

	h.logger.Infof("Review by userID [%v] to productID [%v] was deleted", token.UserID, productID)

	product, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviewQuery := model.ReviewQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	product.Reviews, err = h.services.Review.GetAll(productID, reviewQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productQuery := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     limitRelatedProducts,
			SortBy:    model.SortByViews,
			SortOrder: model.DESCENDING,
		},
		ProductID: productID,
	}

	product.RelatedProducts, err = h.services.Product.GetProductsByCategory(product.Category, productQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(product); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}
