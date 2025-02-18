package user

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/entity"
	"net/http"
)

type httpHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *httpHandler {
	return &httpHandler{
		service: service,
	}
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Create(r.Context(), user); err != nil {
		resp.WriteError(w, http.StatusInternalServerError, err)
	}

	resp.WriteJSON(w, http.StatusCreated, user)
}
