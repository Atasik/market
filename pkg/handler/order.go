package handler

import (
	"encoding/json"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	session, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	order := model.Order{
		CreatedAt:   time.Now(),
		DeliveredAt: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.Services.Order.Create(session.ID, order)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Insert into Orders with id LastInsertId: %v", lastID)

	orders, err := h.Services.Order.GetAll(session.ID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusCreated)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	session, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	selectedOrder, err := h.Services.Order.GetByID(session.ID, orderID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedOrder.Products, err = h.Services.Order.GetProductsByOrderID(session.ID, orderID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(selectedOrder)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}
