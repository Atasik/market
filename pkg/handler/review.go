package handler

import (
	"encoding/json"
	"market/pkg/model"
	"market/pkg/session"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadGateway)
	}

	r.ParseForm()
	review := model.Review{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&review, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
	}

	review.CreationDate = time.Now()
	review.ProductID = productID
	review.UserID = sess.UserID
	review.Username = sess.UserName

	review.ID, err = h.Repository.ReviewRepo.Create(review)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	_, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadGateway)
	}

	_, err = h.Repository.ReviewRepo.Delete(reviewID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, err := json.Marshal(map[string]uint32{
		"updated": uint32(reviewID),
	})
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadGateway)
	}

	r.ParseForm()
	review := model.Review{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&review, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
	}

	_, err = h.Repository.ReviewRepo.Update(sess.UserID, productID, review.Text)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}
