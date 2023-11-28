package v1

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
	resp, _ := json.Marshal(errorResponse{msg}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newStatusReponse(w http.ResponseWriter, msg string, status int) {
	resp, _ := json.Marshal(statusResponse{msg}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newGetProductsResponse(w http.ResponseWriter, products []model.Product, status int) {
	resp, _ := json.Marshal(getProductsResponse{products}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}

func newGetOrdersResponse(w http.ResponseWriter, orders []model.Order, status int) {
	resp, _ := json.Marshal(getOrdersResponse{orders}) //nolint:errcheck
	w.WriteHeader(status)
	w.Write(resp) //nolint:errcheck
}
