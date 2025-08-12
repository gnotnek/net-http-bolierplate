package category

import (
	"encoding/json"
	"net-http-boilerplate/internal/api/resp"
	"net-http-boilerplate/internal/entity"
	apperror "net-http-boilerplate/internal/pkg/app-error"
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
		log.Ctx(ctx).Error().Err(err).Msg("failed to decode request")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	if err := h.service.Create(ctx, &req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("failed to create category: %s", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", nil)

}

func (h *httpHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("perPage")

	if pageStr == "" {
		pageStr = "1"
	}
	if perPageStr == "" {
		perPageStr = "10"
	}

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("invalid 'page' query param")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "'page' must be a number"))
		return
	}

	perPageInt, err := strconv.Atoi(perPageStr)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("invalid 'perPage' query param")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "'perPage' must be a number"))
		return
	}

	filter := &entity.Filter{
		Page:    &pageInt,
		PerPage: &perPageInt,
	}

	res, stats, err := h.service.FindAll(ctx, filter)
	if err != nil {
		if err == apperror.ErrResourceNotFound {
			resp.WriteJSONWithPaginateResponse(w, http.StatusOK, "success", res, stats)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to fetch categories: %v", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSONWithPaginateResponse(w, http.StatusOK, "success", res, stats)
}

func (h *httpHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("bad request, invalid id")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	category, err := h.service.FindByID(ctx, id)
	if err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("category not found")
			resp.WriteError(w, resp.NewError(http.StatusNotFound, "category not found"))
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to get category: %s", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", category)
}

func (h *httpHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("bad request, invalid id")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("bad request")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	res, err := h.service.Update(ctx, id, req)
	if err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("category not found")
			resp.WriteError(w, resp.NewError(http.StatusNotFound, "category not found"))
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to update category: %s", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", res)
}

func (h *httpHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("bad request, invalid id")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "bad request"))
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("category not found")
			resp.WriteError(w, resp.NewError(http.StatusNotFound, "category not found"))
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("failed to delete category: %s", err)
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", nil)
}
