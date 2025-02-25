package category

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/entity"
	"net/http"
	"strconv"

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

func (h *httpHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.service.FindAll(ctx)
	if err != nil {
		if err == ErrCategoryNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("categories not found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "categories not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to get categories")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, categories)
}

func (h *httpHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid request")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	category, err := h.service.FindByID(ctx, id)
	if err != nil {
		if err == ErrCategoryNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("category not found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "category not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to get category")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, category)
}

func (h *httpHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid request")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	category := &entity.Category{
		ID:   id,
		Name: req.Name,
	}

	if err := h.service.Update(ctx, category); err != nil {
		if err == ErrCategoryNotFound {
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "category not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to update category")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]string{"message": "category updated"})
}

func (h *httpHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid request")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if err == ErrCategoryNotFound {
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "category not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to delete category")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
}
