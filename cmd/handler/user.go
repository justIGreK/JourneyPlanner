package handler

import (
	"JourneyPlanner/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// @Summary SignUp
// @Tags users
// @Description Create account
// @Produce  json
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Param email query string true "your email"
// @Router /auth/singUp [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	input := models.SignUp{
		Login:    r.URL.Query().Get("login"),
		Email:    r.URL.Query().Get("email"),
		Password: r.URL.Query().Get("password"),
	}
	if err := validate.Struct(input); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.User.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary SignIn
// @Tags users
// @Description Authorization to the account
// @Produce  json
// @Param option query string true "your login or email"
// @Param password query string true "your password"
// @Router /auth/signIn [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	credentials := models.LoginRequest{
		Option:   r.URL.Query().Get("option"),
		Password: r.URL.Query().Get("password"),
	}
	if err := validate.Struct(credentials); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.User.LoginUser(r.Context(), credentials.Option, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err =json.NewEncoder(w).Encode(token)
	if err != nil {
		logs.Error("failed to encode JSON: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
