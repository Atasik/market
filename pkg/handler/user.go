package handler

import (
	"market/pkg/model"
	"net/http"

	"github.com/gorilla/schema"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := model.User{}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := decoder.Decode(&user, r.PostForm)
	if err != nil {
		http.Error(w, `Bad form`, http.StatusBadRequest)
	}

	id, err := h.Services.CreateUser(user)
	if err != nil {
		http.Error(w, `Failed to verify password`, http.StatusInternalServerError)
		return
	}

	sess, err := h.Sessions.Create(w, id, user.Username, "user")
	if err != nil {
		http.Error(w, `Session Error`, http.StatusUnauthorized)
		return
	}

	_, err = h.Services.Basket.CreateBasket(id)
	if err != nil {
		http.Error(w, "Create Basket Error", http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("username")
	password := r.FormValue("password")

	u, err := h.Services.User.VerifyUser(login, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sess, err := h.Sessions.Create(w, u.ID, u.Username, u.Role)
	if err != nil {
		http.Error(w, `Session Error`, http.StatusUnauthorized)
		return
	}
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
