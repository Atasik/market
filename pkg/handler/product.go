package handler

import (
	"encoding/json"
	"fmt"
	"market/pkg/model"
	"market/pkg/session"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	restrictedMsg = "Access denied, you are not admin"
)

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	sess, err := session.SessionFromContext(r.Context())
	if err == nil {
		products, err := h.Repository.BasketRepo.GetByID(sess.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, prd := range products {
			sess.AddPurchase(prd.ID)
		}
	}

	products, err := h.Repository.ProductRepo.GetAll(orderBy)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Products   []model.Product
		Session    *session.Session
		TotalCount int
	}{
		Products:   products,
		Session:    sess,
		TotalCount: 0,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Product(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		print("no sess")
	}
	//bsk := basket.Basket{}
	//if err == nil {
	//bsk, err = h.BasketRepo.GetByID(sess.UserID)
	// if err != nil {
	// 	http.Error(w, `Database Error`, http.StatusInternalServerError)
	// 	return
	// }

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	selectedProduct, err := h.Repository.ProductRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	input := model.UpdateProductInput{
		Views: &selectedProduct.Views,
	}

	_, err = h.Repository.ProductRepo.Update(selectedProduct.ID, input)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	relatedProducts, err := h.Repository.ProductRepo.GetByType(selectedProduct.Type, 5)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "product.html", struct {
		Product    model.Product
		Related    []model.Product
		Session    *session.Session
		TotalCount int
	}{
		Product:    selectedProduct,
		Related:    relatedProducts,
		Session:    sess,
		TotalCount: 0,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	if sess.UserType != "admin" {
		http.Error(w, restrictedMsg, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad Id", http.StatusBadGateway)
		return
	}

	_, err = h.Repository.ProductRepo.Delete(id)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	sess.DeletePurchase(id)
	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]uint32{
		"updated": uint32(id),
	})
	w.Write(respJSON)
}

func (h *Handler) AddProductForm(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	if sess.UserType != "admin" {
		http.Error(w, restrictedMsg, http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "create_product.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	if sess.UserType != "admin" {
		http.Error(w, restrictedMsg, http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(10 << 20)
	product := model.Product{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&product, r.PostForm)
	if err != nil {
		print(err.Error())
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}

	url, err := h.ImageService.Upload(file)
	if err != nil {
		http.Error(w, `ImageService Error`, http.StatusInternalServerError)
		return
	}

	product.ImageURL = url

	defer file.Close()

	lastID, err := h.Repository.ProductRepo.Create(product)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert into Products with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) UpdateProductForm(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	if sess.UserType != "admin" {
		http.Error(w, restrictedMsg, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	prod, err := h.Repository.ProductRepo.GetByID(id)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "update_product.html", struct {
		Product    model.Product
		Session    *session.Session
		TotalCount int
	}{
		Product:    prod,
		Session:    sess,
		TotalCount: 0,
	})
	if err != nil {
		print(err.Error())
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}
	if sess.UserType != "admin" {
		http.Error(w, restrictedMsg, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var input model.UpdateProductInput
	err = decoder.Decode(&input, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	switch err {
	case nil:
		url, err := h.ImageService.Upload(file)
		if err != nil {
			http.Error(w, `ImageService Error`, http.StatusInternalServerError)
			return
		}

		input.ImageURL = &url
		defer file.Close()
	case http.ErrMissingFile:
		fmt.Println("no file")
	default:
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}

	ok, err := h.Repository.ProductRepo.Update(id, input)
	if err != nil {
		http.Error(w, `Database error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", "heh", ok)

	http.Redirect(w, r, "/", http.StatusFound)
}
