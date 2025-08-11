package user

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	apperror "net-http-boilerplate/internal/pkg/app-error"
	"net-http-boilerplate/internal/pkg/validator"
	"net/http"

	"github.com/rs/zerolog/log"
)

type httpHandler struct {
	service   *Service
	validator *validator.Validator
}

func NewUserHandler(service *Service, validator *validator.Validator) *httpHandler {
	return &httpHandler{
		service:   service,
		validator: validator,
	}
}

func (h *httpHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("cannot decode request")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	if err := h.service.Register(ctx, &req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("cannot register new user: %s", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusCreated, "success", nil)
}

func (h *httpHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteError(w, err)
		return
	}

	data, err := h.service.Login(ctx, req)
	switch err {
	case nil:
	case apperror.ErrResourceNotFound:
		resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "user not found"})
		return
	case apperror.ErrInvalidPassword:
		resp.WriteJSON(w, http.StatusUnauthorized, map[string]string{"message": "invalid password"})
		return
	default:
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", data)
}
