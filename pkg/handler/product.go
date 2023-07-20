package handler

import (
	"context"
	"encoding/json"
	"market/pkg/model"
	"market/pkg/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const imageUploadTimeout = 5 * time.Second

func (h *Handler) about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	resp, err := json.Marshal(map[string]interface{}{
		"message": "a simple market-api",
	})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	orderBy := r.URL.Query().Get("order_by")

	products, err := h.Services.Product.GetAll(orderBy)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	selectedProduct, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.Services.Product.IncreaseViewsCounter(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedProduct.Reviews, err = h.Services.Review.GetAll(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedProduct.RelatedProducts, err = h.Services.Product.GetByType(selectedProduct.Category, productID, 5)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(selectedProduct)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	session, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	r.ParseMultipartForm(10 << 20)
	product := model.Product{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = decoder.Decode(&product, r.PostForm)
	if err != nil {
		newErrorResponse(w, `Bad form`, http.StatusBadRequest)
		return
	}

	err = h.Validator.Struct(product)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		newErrorResponse(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()
	data, err := h.Services.Image.Upload(ctx, file)
	if err != nil {
		newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
		return
	}

	product.UserID = session.ID
	product.ImageURL = data.ImageURL
	product.ImageID = data.ImageID
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	defer file.Close()

	lastID, err := h.Services.Product.Create(product)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()
		err = h.Services.Image.Delete(ctx, product.ImageID)
		if err != nil {
			newErrorResponse(w, `ImageServer Error`, http.StatusInternalServerError)
			return
		}
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.ID = lastID

	h.Logger.Infof("Product was created with id LastInsertId: %v", lastID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, `Bad id`, http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(10 << 20)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var input model.UpdateProductInput
	err = decoder.Decode(&input, r.PostForm)
	if err != nil {
		newErrorResponse(w, `Bad form`, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	noFile := err == http.ErrMissingFile
	print(noFile)
	if err != nil && !noFile {
		newErrorResponse(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}

	if !noFile {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()
		data, err := h.Services.Image.Upload(ctx, file)
		if err != nil {
			newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
			return
		}
		input.ImageURL = &data.ImageURL
		input.ImageID = &data.ImageID
		defer file.Close()
	}

	currentTime := time.Now()
	input.UpdatedAt = &currentTime

	oldProduct, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := h.Services.Product.Update(sess.ID, productID, input)
	if err != nil {
		if !noFile {
			ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
			defer cancel()
			err = h.Services.Image.Delete(ctx, *input.ImageID)
			if err != nil {
				newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
				return
			}
		}
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !noFile {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()
		err = h.Services.Image.Delete(ctx, oldProduct.ImageID)
		if err != nil {
			newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
			return
		}
	}

	product, err := h.Services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Product was updated: %v %v", product, ok)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	session, err := service.SessionFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Session Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	product, err := h.Services.Product.GetByID(productId)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := h.Services.Product.Delete(session.ID, productId)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Product was deleted: %v %v", product, ok)

	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()
	err = h.Services.Image.Delete(ctx, product.ImageID)
	if err != nil {
		newErrorResponse(w, "ImageService Error", http.StatusInternalServerError)
		return
	}

	newStatusReponse(w, "done", http.StatusOK)
}
