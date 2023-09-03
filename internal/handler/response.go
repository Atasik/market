package handler

import (
	"encoding/json"
	"market/internal/model"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

type getProductsResponse struct {
	Data []model.Product `json:"data"`
}

type getOrdersResponse struct {
	Data []model.Order `json:"data"`
}

func newErrorResponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(errorResponse{msg})
	w.WriteHeader(status)
	w.Write(resp)
}

func newStatusReponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(statusResponse{msg})
	w.WriteHeader(status)
	w.Write(resp)
}

func newGetProductsResponse(w http.ResponseWriter, products []model.Product, status int) {
	resp, _ := json.Marshal(getProductsResponse{products})
	w.WriteHeader(status)
	w.Write(resp)
}

func newGetOrdersResponse(w http.ResponseWriter, orders []model.Order, status int) {
	resp, _ := json.Marshal(getOrdersResponse{orders})
	w.WriteHeader(status)
	w.Write(resp)
}
