package user

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/entity"
	"net-http-boilerplate/internal/jwt"
	"net/http"
)

type httpHandler struct {
	service *Service
	jwt     *jwt.JWT
}

func NewHTTPHandler(service *Service, jwt *jwt.JWT) *httpHandler {
	return &httpHandler{
		service: service,
		jwt:     jwt,
	}
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteError(w, err)
		return
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Register(r.Context(), user); err != nil {
		resp.WriteError(w, err)
	}

	resp.WriteJSON(w, http.StatusCreated, user)
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteError(w, err)
		return
	}

	user, err := h.service.Login(ctx, req)
	switch err {
	case nil:
	case ErrUserNotFound:
		resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "user not found"})
		return
	case ErrInvalidPassword:
		resp.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid password"})
		return
	default:
		resp.WriteError(w, err)
		return
	}

	accessToken, refreshToken, err := h.jwt.GenerateToken(user.ID, user.Email)
	if err != nil {
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"data":          user,
	})
}

func (h *httpHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	newToken, err := h.jwt.RefreshToken(tokenString)
	if err != nil {
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]string{"access_token": newToken})
}
