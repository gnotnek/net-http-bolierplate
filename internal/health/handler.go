package health

import (
	"net-http-boilerplate/internal/api/resp"
	"net/http"
)

type httpHandler struct {
	service *Service
}

func NewHttpHandler(service *Service) *httpHandler {
	return &httpHandler{
		service: service,
	}
}

func (h *httpHandler) CheckDatabase(w http.ResponseWriter, r *http.Request) {
	healthComponent, isHealthy := h.service.CheckDatabase(r.Context())
	statusCode := http.StatusOK
	if !isHealthy {
		statusCode = http.StatusInternalServerError
	}

	resp.WriteJSON(w, statusCode, healthComponent)
}
