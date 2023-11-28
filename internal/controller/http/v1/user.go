package v1

import (
	"encoding/json"
	"io"
	"market/internal/model"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) initUserRoutes(api *mux.Router) {
	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/sign-in", h.signIn).Methods("POST")
	user.HandleFunc("/sign-up", h.signUp).Methods("POST")
	user.HandleFunc("/{userId}/products", queryMiddleware(h.getProductsByUserID)).Methods("GET")
}

type signUpInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary	Register in the market
// @Tags		user
// @ID			register
// @Accept		json
// @Produce	json
// @Param		input	body		model.User	true	"Account info"
// @Success	200		{string}	string		"token"
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/user/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var inp signUpInput
	if err = json.Unmarshal(body, &inp); err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	if err = h.validator.Struct(inp); err != nil {
		newErrorResponse(w, "invalid input", http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: inp.Username,
		Password: inp.Password,
	}
	userID, err := h.services.User.CreateUser(user)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = h.services.Cart.Create(userID); err != nil {
		newErrorResponse(w, "Create Basket Error", http.StatusInternalServerError)
		return
	}

	token, err := h.services.User.GenerateToken(user.Username, user.Password)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}

type signInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary	Login into market
// @Tags		user
// @ID			login
// @Accept		json
// @Produce	json
// @Param		input	body		signInInput	true	"Username and password"
// @Success	200		{string}	string		"token"
// @Failure	400,404	{object}	errorResponse
// @Failure	500		{object}	errorResponse
// @Failure	default	{object}	errorResponse
// @Router		/api/v1/user/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", appJSON)
	if r.Header.Get("Content-Type") != appJSON {
		newErrorResponse(w, "unknown payload", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, "server error", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	var inp signInInput
	if err = json.Unmarshal(body, &inp); err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	if err = h.validator.Struct(inp); err != nil {
		newErrorResponse(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.services.User.GenerateToken(inp.Username, inp.Password)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})
	if err != nil {
		newErrorResponse(w, `can't create payload`, http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(resp); err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}
