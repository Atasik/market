package handler

import (
	"encoding/json"
	"io/ioutil"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
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

	review.CreatedAt = time.Now()
	review.ProductID = productID
	review.UserID = sess.ID
	review.Username = sess.Username

	review.ID, err = h.Services.Review.Create(review)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviews, err := h.Services.Review.GetAll(productID, "")
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatedProducts, err := h.Services.Product.GetByType(product.Category, 5)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews = reviews
	product.RelatedProducts = relatedProducts

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	session, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["reviewId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
	}

	_, err = h.Services.Review.Delete(session.ID, reviewID)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad id", http.StatusBadRequest)
	}

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	reviews, err := h.Services.Review.GetAll(productID, "")
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatedProducts, err := h.Services.Product.GetByType(product.Category, 5)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews = reviews
	product.RelatedProducts = relatedProducts

	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
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
	}

	body, err := ioutil.ReadAll(r.Body)
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

	currentTime := time.Now()
	input.UpdatedAt = &currentTime

	_, err = h.Services.Review.Update(sess.ID, productID, input)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	reviews, err := h.Services.Review.GetAll(productID, "")
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	relatedProducts, err := h.Services.Product.GetByType(product.Category, 5)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Reviews = reviews
	product.RelatedProducts = relatedProducts

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}
