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

const (
	imageUploadTimeout   = 5 * time.Second
	limitRelatedProducts = 5
)

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	sess, err := service.SessionFromContext(r.Context())
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

	product.UserID = sess.UserID
	product.ImageURL = data.ImageURL
	product.ImageID = data.ImageID

	defer file.Close()

	productID, err := h.Services.Product.Create(product)
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

	product.ID = productID

	h.Logger.Infof("Product was created with id LastInsertId: %v", productID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.ProductQueryInput{
		Limit:     options.Limit,
		Offset:    options.Offset,
		SortBy:    options.SortBy,
		SortOrder: options.SortOrder,
	}

	err = q.Validate()
	if err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.Services.Product.GetAll(q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) getProductsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.ProductQueryInput{
		Limit:     options.Limit,
		Offset:    options.Offset,
		SortBy:    options.SortBy,
		SortOrder: options.SortOrder,
	}

	err = q.Validate()
	if err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.Services.Product.GetProductsByUserID(userID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) getProductsByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	vars := mux.Vars(r)
	categoryName := vars["categoryName"]

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.ProductQueryInput{
		Limit:     options.Limit,
		Offset:    options.Offset,
		SortBy:    options.SortBy,
		SortOrder: options.SortOrder,
	}

	err = q.Validate()
	if err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.Services.Product.GetProductsByCategory(categoryName, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

func (h *Handler) getProductByID(w http.ResponseWriter, r *http.Request) {
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

	err = h.Services.Product.IncreaseViewsCounter(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectedProduct.Reviews, err = h.Services.Review.GetAll(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := model.ProductQueryInput{
		Limit:     5,
		Offset:    0,
		ProductID: productID,
		SortBy:    model.SortByViews,
		SortOrder: model.DESCENDING,
	}

	selectedProduct.RelatedProducts, err = h.Services.Product.GetProductsByCategory(selectedProduct.Category, q)
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

	err = h.Services.Product.Update(sess.UserID, productID, input)
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

	h.Logger.Infof("Product was updated: %v", product)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	sess, err := service.SessionFromContext(r.Context())
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

	err = h.Services.Product.Delete(sess.UserID, productId)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("Product was deleted: %v", product)

	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()
	err = h.Services.Image.Delete(ctx, product.ImageID)
	if err != nil {
		newErrorResponse(w, "ImageService Error", http.StatusInternalServerError)
		return
	}

	newStatusReponse(w, "done", http.StatusOK)
}
