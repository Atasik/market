package handler

import (
	"encoding/json"
	"fmt"
	"market/pkg/model"
	"market/pkg/session"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) AddProductToBasket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Database Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	basketId, err := h.Repository.BasketRepo.AddProduct(sess.UserID, productId)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	sess.AddPurchase(productId)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]int{
		"updated": basketId,
	})
	w.Write(respJSON)
}

func (h *Handler) DeleteProductFromBasket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	_, err = h.Repository.BasketRepo.DeleteProduct(sess.UserID, int(id))
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	sess.DeletePurchase(int(id))
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]uint32{
		"updated": uint32(id),
	})
	w.Write(respJSON)
}

func (h *Handler) Basket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	products := []model.Product{}
	if err == nil {
		products, err = h.Repository.BasketRepo.GetByID(sess.UserID)
		if err != nil {
			http.Error(w, `Database Error`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "basket.html", struct {
		Products   []model.Product
		TotalPrice int
		Session    *session.Session
		TotalCount int
	}{
		Products:   products,
		TotalPrice: 0,
		Session:    sess,
		TotalCount: 0,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}