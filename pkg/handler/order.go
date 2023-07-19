package handler

import (
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"time"
)

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	orders, err := h.Services.Order.GetAll(sess.ID)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
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

	order := model.Order{
		CreatedAt:   time.Now(),
		DeliveredAt: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.Services.Order.Create(sess.ID, order)
	if err != nil {
		newErrorResponse(w, `Database error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Insert into Orders with id LastInsertId: %v", lastID)

	orders, err := h.Services.Order.GetAll(sess.ID)
	if err != nil {
		newErrorResponse(w, "Database Error", http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusCreated)
}
