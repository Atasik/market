package handler

import (
	"encoding/json"
	"io"
	"market/internal/model"
	"net/http"
)

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
// @Router		/api/register [post]
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

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	user.Role = model.USER

	err = h.validator.Struct(user)
	if err != nil {
		newErrorResponse(w, "invalid input", http.StatusBadRequest)
		return
	}

	userID, err := h.services.User.CreateUser(user)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.services.Cart.Create(userID)
	if err != nil {
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

	_, err = w.Write(resp)
	if err != nil {
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
// @Router		/api/login [post]
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

	var input signInInput
	err = json.Unmarshal(body, &input)
	if err != nil {
		newErrorResponse(w, "cant unpack payload", http.StatusBadRequest)
		return
	}

	err = h.validator.Struct(input)
	if err != nil {
		newErrorResponse(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.services.User.GenerateToken(input.Username, input.Password)
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

	_, err = w.Write(resp)
	if err != nil {
		newErrorResponse(w, `can't write resp`, http.StatusInternalServerError)
		return
	}
}
