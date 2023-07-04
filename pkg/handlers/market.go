package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"market/pkg/basket"
	"market/pkg/order"
	"market/pkg/product"
	"market/pkg/services"
	"market/pkg/session"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
)

type MarketHandler struct {
	Tmpl         *template.Template
	Logger       *zap.SugaredLogger
	Sessions     *session.SessionsManager
	ProductRepo  product.ProductRepo
	OrderRepo    order.OrderRepo
	BasketRepo   basket.BasketRepo
	ImageService services.ImageService
}

const (
	restrictedMsg = "Access denied, you are not admin"
)

func (h *MarketHandler) Index(w http.ResponseWriter, r *http.Request) {
	orderBy := r.URL.Query().Get("order_by")
	sess, err := session.SessionFromContext(r.Context())
	if err == nil {
		products, err := h.BasketRepo.GetByID(sess.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, prd := range products {
			sess.AddPurchase(prd.ID)
		}
	}

	products, err := h.ProductRepo.GetAll(orderBy)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Products   []product.Product
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

func (h *MarketHandler) About(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *MarketHandler) Privacy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "privacy.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *MarketHandler) Product(w http.ResponseWriter, r *http.Request) {
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

	selectedProduct, err := h.ProductRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	input := product.UpdateProductInput{
		Views: &selectedProduct.Views,
	}

	_, err = h.ProductRepo.Update(selectedProduct.ID, input)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	relatedProducts, err := h.ProductRepo.GetByType(selectedProduct.Type, 5)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "product.html", struct {
		Product    product.Product
		Related    []product.Product
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

func (h *MarketHandler) AddProductToBasket(w http.ResponseWriter, r *http.Request) {
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

	basketId, err := h.BasketRepo.AddProduct(sess.UserID, productId)
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

func (h *MarketHandler) DeleteProductFromBasket(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.BasketRepo.DeleteProduct(sess.UserID, int(id))
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

func (h *MarketHandler) History(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}

	orders, err := h.OrderRepo.GetAll(sess.UserID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "history.html", struct {
		Landings []order.Order
	}{
		Landings: orders,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MarketHandler) Basket(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	products := []product.Product{}
	if err == nil {
		products, err = h.BasketRepo.GetByID(sess.UserID)
		if err != nil {
			http.Error(w, `Database Error`, http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "basket.html", struct {
		Products   []product.Product
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

func (h *MarketHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.ProductRepo.Delete(id)
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

func (h *MarketHandler) AddProductForm(w http.ResponseWriter, r *http.Request) {
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

func (h *MarketHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
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
	product := product.Product{}
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

	lastID, err := h.ProductRepo.Create(product)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Insert into Products with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MarketHandler) UpdateProductForm(w http.ResponseWriter, r *http.Request) {
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

	prod, err := h.ProductRepo.GetByID(id)
	if err != nil {
		http.Error(w, `Database Error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = h.Tmpl.ExecuteTemplate(w, "update_product.html", struct {
		Product    product.Product
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

func (h *MarketHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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
	var input product.UpdateProductInput
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

	ok, err := h.ProductRepo.Update(id, input)
	if err != nil {
		http.Error(w, `Database error`, http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("update: %v %v", "heh", ok)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *MarketHandler) RegisterOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		http.Error(w, "Session Error", http.StatusBadRequest)
		return
	}

	products, err := h.BasketRepo.GetByID(sess.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := order.Order{
		CreationDate: time.Now(),
		DeliveryDate: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.OrderRepo.Create(sess.UserID, order, products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Logger.Infof("Insert into Orders with id LastInsertId: %v", lastID)
	http.Redirect(w, r, "/", http.StatusFound)
}
