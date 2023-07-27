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

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	order := model.Order{
		DeliveredAt: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.Services.Order.Create(sess.UserID, order)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Order was created with id LastInsertId: %v", lastID)

	orders, err := h.Services.Order.GetAll(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusCreated)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
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

	selectedOrder, err := h.Services.Order.GetByID(sess.UserID, orderID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedOrder.Products, err = h.Services.Order.GetProductsByOrderID(sess.UserID, orderID)
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

func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	orders, err := h.Services.Order.GetAll(sess.UserID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusOK)
}
