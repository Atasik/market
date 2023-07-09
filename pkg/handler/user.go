package handler

import (
	"market/pkg/repository"
	"net/http"
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
	login := r.FormValue("login")
	password := r.FormValue("password")

	id, err := h.Repository.Register(login, password)
	if err == repository.ErrUserExists {
		http.Error(w, `User already exists Error`, http.StatusBadRequest)
		return
	}

	sess, _ := h.Sessions.Create(w, id, login, "user")
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

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	u, err := h.Repository.UserRepo.Authorize(r.FormValue("login"), r.FormValue("password"))
	if err == repository.ErrNoUser {
		http.Error(w, `User doesn't exist Error`, http.StatusUnauthorized)
		return
	}
	if err == repository.ErrBadPass {
		http.Error(w, `Bad password Error`, http.StatusUnauthorized)
		return
	}

	sess, _ := h.Sessions.Create(w, u.ID, u.Username, u.UserMode)
	h.Logger.Infof("created session for %v", sess.UserID)
	http.Redirect(w, r, "/", http.StatusFound)
}
