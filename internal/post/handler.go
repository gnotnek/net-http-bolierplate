package post

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

	post := &entity.Post{
		Title:      req.Title,
		Content:    req.Content,
		ImageURL:   req.ImageURL,
		CategoryID: req.CategoryID,
	}

	if err := h.service.Create(ctx, post); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to create post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusCreated, post)
}

func (h *httpHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := h.service.FindAll(ctx)
	if err != nil {
		if err == ErrPostNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("no posts found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "post not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch posts")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, posts)
}

func (h *httpHandler) FindByCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	category := r.URL.Query().Get("category")
	if category == "" {
		log.Ctx(ctx).Error().Msg("category is required")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "category is required"})
		return
	}

	posts, err := h.service.FindByCategory(ctx, category)
	if err != nil {
		if err == ErrPostNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("no posts found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "no posts found"})
			return
		}

		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, posts)
}

func (h *httpHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid id")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	post, err := h.service.FindByID(ctx, id)
	if err != nil {
		if err == ErrPostNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("post not found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "post not found"})
			return
		}

		log.Ctx(ctx).Error().Err(err).Msg("failed to fetch post")
		resp.WriteError(w, err)
		return
	}

	resp.WriteJSON(w, http.StatusOK, post)
}

func (h *httpHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("invalid id")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	var req UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to decode request")
		resp.WriteJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request"})
		return
	}

	post := &entity.Post{
		ID:         id,
		Title:      req.Title,
		Content:    req.Content,
		ImageURL:   req.ImageURL,
		CategoryID: req.CategoryID,
	}

	if err := h.service.Update(ctx, post); err != nil {
		if err == ErrPostNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("post not found")
			resp.WriteJSON(w, http.StatusNotFound, map[string]string{"message": "post not found"})
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
		if err == ErrPostNotFound {
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
