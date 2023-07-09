package handler

import (
	"market/pkg/model"
	"market/pkg/session"
	"net/http"
	"time"
)

func (h *Handler) History(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}

	orders, err := h.Repository.OrderRepo.GetAll(sess.UserID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "history.html", struct {
		Landings   []model.Order
		Session    *session.Session
		TotalCount int
	}{
		Landings:   orders,
		Session:    sess,
		TotalCount: 0,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RegisterOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}

	products, err := h.Repository.BasketRepo.GetByID(sess.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := model.Order{
		CreationDate: time.Now(),
		DeliveryDate: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.Repository.OrderRepo.Create(sess.UserID, order, products)
	if err != nil {
		http.Error(w, `Database error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Insert into Orders with id LastInsertId: %v", lastID)

	_, err = h.Repository.BasketRepo.DeleteAll(sess.UserID)
	if err != nil {
		http.Error(w, `Database error`, http.StatusInternalServerError)
		return
	}
	sess.PurgeBasket()

	http.Redirect(w, r, "/", http.StatusFound)
}
