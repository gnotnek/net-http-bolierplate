package post

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

func NewPostHandler(service *Service) *httpHandler {
	return &httpHandler{
		service: service,
	}
}

func (h *httpHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to decode request")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	data, err := h.service.Create(ctx, &req)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", data)
}

func (h *httpHandler) FindAll(w http.ResponseWriter, r *http.Request) {
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

	posts, stats, err := h.service.FindAll(ctx, filter)
	if err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("no posts found")
			resp.WriteJSONWithPaginateResponse(w, http.StatusOK, "success", posts, stats)
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch posts")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSONWithPaginateResponse(w, http.StatusOK, "success", posts, stats)
}

func (h *httpHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid id")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "invalid id"))
		return
	}

	post, err := h.service.FindByID(ctx, id)
	if err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("post not found")
			resp.WriteError(w, resp.NewError(http.StatusNotFound, "post not found"))
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteSuccess(w, http.StatusOK, "success", post)
}

func (h *httpHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid id")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "invalid id"))
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to decode request")
		resp.WriteError(w, resp.NewError(http.StatusBadRequest, "invalid request"))
		return
	}

	post := &entity.Post{
		ID:         id,
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
	}

	if err := h.service.Update(ctx, post); err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("post not found")
			resp.WriteError(w, resp.NewError(http.StatusBadRequest, "post not found"))
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to update post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, post)
}

func (h *httpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid id")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		if err == apperror.ErrResourceNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("post not found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "post not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to delete post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, map[string]string{"message": "post deleted"})
}
