package v1

import (
	"context"
	"encoding/json"
	"market/internal/model"
	"market/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	imageUploadTimeout   = 5 * time.Second
	limitRelatedProducts = 5
	limitFileBytes       = 10 << 20
)

func (h *Handler) initProductRoutes(api *mux.Router) {
	product := api.PathPrefix("/product").Subrouter()
	product.HandleFunc("", h.authMiddleware(h.createProduct)).Methods("POST")
	product.HandleFunc("/{productId}", queryMiddleware(h.getProductByID)).Methods("GET")
	product.HandleFunc("/{productId}", h.authMiddleware(h.updateProduct)).Methods("PUT")
	product.HandleFunc("/{productId}", h.authMiddleware(h.deleteProduct)).Methods("DELETE")
	product.HandleFunc("/{productId}/review", h.authMiddleware(h.createReview)).Methods("POST")
	product.HandleFunc("/{productId}/review/{reviewId}", h.authMiddleware(h.updateReview)).Methods("PUT")
	product.HandleFunc("/{productId}/review/{reviewId}", h.authMiddleware(h.deleteReview)).Methods("DELETE")
}

func (h *Handler) initProductsRoutes(api *mux.Router) {
	products := api.PathPrefix("/products").Subrouter()
	products.Methods("GET").HandlerFunc(queryMiddleware(h.getAllProducts))
	products.HandleFunc("/category/{categoryName}", queryMiddleware(h.getProductsByCategory)).Methods("GET")
}

// @Summary	Add a new product to the market
// @Security	ApiKeyAuth
// @Tags		products
// @ID			create-product
// @Accept		mpfd
// @Product	json
// @Param		file		formData	file	true	"Image to Upload"
// @Param		title		formData	string	true	"Title of product"
// @Param		price		formData	number	true	"Price of product"
// @Param		tag			formData	string	false	"Tag of product"
// @Param		category	formData	string	true	"Category of product"
// @Param		description	formData	string	false	"Description of product"
// @Param		amount		formData	integer	true	"Amount of products"
// @Success	201			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/v1/product [post]
func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	if err = r.ParseMultipartForm(limitFileBytes); err != nil {
		newErrorResponse(w, "Failed to Parse MultipartForm", http.StatusInternalServerError)
		return
	}
	var product model.Product
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err = decoder.Decode(&product, r.PostForm); err != nil {
		newErrorResponse(w, `Bad form`, http.StatusBadRequest)
		return
	}

	if err = h.validator.Struct(product); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		newErrorResponse(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()
	data, err := h.services.Image.Upload(ctx, file)
	if err != nil {
		newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
		return
	}

	product.UserID = token.UserID
	product.ImageURL = data.ImageURL
	product.ImageID = data.ImageID
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	productID, err := h.services.Product.Create(product)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()

		if err = h.services.Image.Delete(ctx, product.ImageID); err != nil {
			newErrorResponse(w, `ImageServer Error`, http.StatusInternalServerError)
			return
		}
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.ID = productID

	h.logger.Info("Product was created", map[string]interface{}{"lastInsetedId": productID})

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(product); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Get all products from the market
// @Tags		products
// @ID			get-all-products
// @Product	json
// @Param   sort_by query   string false "sort by" Enums(views, price, created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		{object}	getProductsResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/products [get]
func (h *Handler) getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	if err = q.Validate(); err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.services.Product.GetAll(q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary	Get products by UserID
// @Tags		products
// @ID			get-products-by-userId
// @Product	json
// @Param		userId	path		integer	true	"ID of user"
// @Param   sort_by query   string false "sort by" Enums(views, price, created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		{object}	getProductsResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/products/{userId} [get]
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
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	if err = q.Validate(); err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.services.Product.GetProductsByUserID(userID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary	Get all products by category from the market
// @Tags		products
// @ID			get-products-by-category
// @Product	json
// @Param		categoryName	path		string	true	"Name of category"
// @Param   sort_by query   string false "sort by" Enums(views, price, created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		{object}	getProductsResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/products/{categoryName} [get]
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
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	if err = q.Validate(); err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	products, err := h.services.Product.GetProductsByCategory(categoryName, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetProductsResponse(w, products, http.StatusOK)
}

// @Summary	Get product by id from the market
// @Tags		products
// @ID			get-product-by-id
// @Product	json
// @Param		productId	path		integer	true	"ID of product to get"
// @Param   sort_by query   string false "sort by" Enums(created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/v1/product/{productId} [get]
func (h *Handler) getProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	selectedProduct, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.services.Product.IncreaseViewsCounter(productID); err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reviewQuery := model.ReviewQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	if err = reviewQuery.Validate(); err != nil {
		newErrorResponse(w, "Bad query", http.StatusBadRequest)
		return
	}

	selectedProduct.Reviews, err = h.services.Review.GetAll(productID, reviewQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productQuery := model.ProductQueryInput{
		QueryInput: model.QueryInput{
			Limit:     limitRelatedProducts,
			SortBy:    model.SortByViews,
			SortOrder: model.DESCENDING,
		},
		ProductID: productID,
	}

	selectedProduct.RelatedProducts, err = h.services.Product.GetProductsByCategory(selectedProduct.Category, productQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(selectedProduct); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Update an existing product from the market
// @Security	ApiKeyAuth
// @Tags		products
// @ID			update-product
// @Accept		mpfd
// @Product	json
// @Param		productId	path		integer	false	"ID of product to update"
// @Param		file		formData	file	false	"Image to Upload"
// @Param		title		formData	string	false	"Title of product"
// @Param		price		formData	number	false	"Price of product"
// @Param		tag			formData	string	false	"Tag of product"
// @Param		category	formData	string	false	"Category of product"
// @Param		description	formData	string	false	"Description of product"
// @Param		amount		formData	integer	false	"Amount of products"
// @Success	200			{object}	model.Product
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/v1/product/{productId} [put]
func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, `Bad id`, http.StatusBadRequest)
		return
	}

	if err = r.ParseMultipartForm(limitFileBytes); err != nil {
		newErrorResponse(w, "Failed to Parse MultipartForm", http.StatusInternalServerError)
		return
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var input model.UpdateProductInput
	if err = decoder.Decode(&input, r.PostForm); err != nil {
		newErrorResponse(w, `Bad form`, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	fileExists := !(err == http.ErrMissingFile)
	if err != nil && fileExists {
		newErrorResponse(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}

	if fileExists {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()
		data, err := h.services.Image.Upload(ctx, file) //nolint:govet
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

	oldProduct, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.services.Product.Update(token.UserID, productID, input); err != nil {
		if fileExists {
			ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
			defer cancel()
			if err = h.services.Image.Delete(ctx, *input.ImageID); err != nil {
				newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
				return
			}
		}
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if fileExists {
		ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
		defer cancel()
		if err = h.services.Image.Delete(ctx, oldProduct.ImageID); err != nil {
			newErrorResponse(w, `ImageService Error`, http.StatusInternalServerError)
			return
		}
	}

	product, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Product was updated", map[string]interface{}{"product": product})

	if err = json.NewEncoder(w).Encode(product); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Delete product from the market
// @Security	ApiKeyAuth
// @Tags		products
// @ID			delete-product
// @Product	json
// @Param		productId	path		integer	true	"ID of product to delete"
// @Success	200			{object}	statusResponse
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/v1/product/{productId} [delete]
func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["productId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	product, err := h.services.Product.GetByID(productID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.services.Product.Delete(token.UserID, productID); err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Product was deleted", map[string]interface{}{"product": product})

	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()
	if err = h.services.Image.Delete(ctx, product.ImageID); err != nil {
		newErrorResponse(w, "ImageService Error", http.StatusInternalServerError)
		return
	}

	newStatusReponse(w, "done", http.StatusOK)
}
