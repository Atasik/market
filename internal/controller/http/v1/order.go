package v1

import (
	"encoding/json"
	"market/internal/model"
	"market/pkg/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (h *Handler) initOrderRoutes(api *mux.Router) {
	order := api.PathPrefix("/order").Subrouter()
	order.HandleFunc("", h.authMiddleware(h.createOrder)).Methods("POST")
	order.HandleFunc("/{orderId}", queryMiddleware(h.authMiddleware(h.getOrder))).Methods("GET")
}

func (h *Handler) initOrdersRoutes(api *mux.Router) {
	orders := api.PathPrefix("/orders").Subrouter()
	orders.Methods("GET").HandlerFunc(queryMiddleware(h.authMiddleware(h.getOrders)))
}

// @Summary	Create order
// @Security	ApiKeyAuth
// @Tags		order
// @ID			create-order
// @Product	json
// @Success	200		{object}	getOrdersResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/order [post]
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	order := model.Order{
		CreatedAt:   time.Now(),
		DeliveredAt: time.Now().Add(4 * 24 * time.Hour),
	}

	lastID, err := h.services.Order.Create(token.UserID, order)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Info("Order was created with id", map[string]interface{}{"lastInsetedId": lastID})

	q := model.OrderQueryInput{
		QueryInput: model.QueryInput{
			Limit:     defaultLimit,
			SortBy:    defaultSortField,
			SortOrder: model.DESCENDING,
		},
	}

	orders, err := h.services.Order.GetAll(token.UserID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusCreated)
}

// @Summary	Get order
// @Security	ApiKeyAuth
// @Tags		order
// @ID			get-order
// @Product	json
// @Param		orderId	path		integer	true	"ID of order to get"
// @Param   sort_by query   string false "sort by" Enums(views, price, created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	201			{object}	model.Order
// @Failure	400,404		{object}	errorResponse
// @Failure	500			{object}	errorResponse
// @Failure	default		{object}	errorResponse
// @Router		/api/v1/order/{orderId} [get]
func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["orderId"])
	if err != nil {
		newErrorResponse(w, "Bad Id", http.StatusBadRequest)
		return
	}

	selectedOrder, err := h.services.Order.GetByID(token.UserID, orderID)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
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

	selectedOrder.Products, err = h.services.Order.GetProductsByOrderID(token.UserID, orderID, q)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(selectedOrder); err != nil {
		newErrorResponse(w, "server error", http.StatusInternalServerError)
		return
	}
}

// @Summary	Get orders
// @Security	ApiKeyAuth
// @Tags		order
// @ID			get-orders
// @Product	json
// @Param   sort_by query   string false "sort by" Enums(created_at)
// @Param   sort_order query string false "sort order" Enums(asc, desc)
// @Param   limit   query int false "limit" Enums(10, 25, 50)
// @Param   page  query int false "page"
// @Success	200		{object}	getOrdersResponse
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/orders [get]
func (h *Handler) getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)

	options, err := optionsFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.TokenFromContext(r.Context())
	if err != nil {
		newErrorResponse(w, "Token Error", http.StatusInternalServerError)
		return
	}

	orderQuery := model.OrderQueryInput{
		QueryInput: model.QueryInput{
			Limit:     options.Limit,
			Offset:    options.Offset,
			SortBy:    options.SortBy,
			SortOrder: options.SortOrder,
		},
	}

	orders, err := h.services.Order.GetAll(token.UserID, orderQuery)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newGetOrdersResponse(w, orders, http.StatusOK)
}
