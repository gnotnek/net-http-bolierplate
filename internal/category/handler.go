package category

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/entity"
	"net/http"

	"github.com/rs/zerolog/log"
)

type httpHandler struct {
	service *Service
}

func NewCategoryHandler(service *Service) *httpHandler {
	return &httpHandler{
		service: service,
	}
}

func (h *httpHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	category := &entity.Category{
		Name: req.Name,
	}

	if err := h.service.Create(ctx, category); err != nil {
		if err == ErrCategoryAlreadyExists {
			resp.WriteJSON(w, http.StatusConflict, map[string]string{"message": "category already exists"})
			return
		}
		if err == ErrCategoryNotFound {
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "category not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to create category")
		resp.WriteError(w, err)
		return
	}
}
