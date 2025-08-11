package resp

import (
	"encoding/json"
	"errors"
	"math"
	"net-http-boilerplate/internal/entity"
	"net/http"
)

// Meta contains pagination details.
type Meta struct {
	Page      int `json:"page"`
	PageTotal int `json:"page_total"`
	Total     int `json:"total"`
}

// DataPaginate wraps paginated data with meta information.
type DataPaginate struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    Meta   `json:"meta"`
}

// HTTPError is the JSON structure for errors.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse is a general response with optional data.
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Empty can be used to indicate an empty object in the response.
type Empty struct{}

type CustomError struct {
	Code    int
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

func (e CustomError) HTTPStatusCode() int {
	return e.Code
}

// Pass spesific error
func NewError(code int, message string) CustomError {
	return CustomError{
		Code:    code,
		Message: message,
	}
}

// WriteJSON writes raw JSON to the client.
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteSuccess sends a successful response with optional data.
func WriteSuccess(w http.ResponseWriter, statusCode int, message string, data any) {
	resp := SuccessResponse{
		Message: message,
	}

	if data != nil {
		resp.Data = data
	}

	WriteJSON(w, statusCode, resp)
}

// WriteJSONWithPaginateResponse sends paginated data with a message.
func WriteJSONWithPaginateResponse(w http.ResponseWriter, statusCode int, message string, data any, stats *entity.Stats) {
	totalPage := int(math.Ceil(float64(stats.Total) / float64(stats.Limit)))
	meta := Meta{
		Page:      stats.Page,
		PageTotal: totalPage,
		Total:     stats.Total,
	}

	response := DataPaginate{
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	WriteJSON(w, statusCode, response)
}

// WriteError sends an error response in a consistent JSON format.
func WriteError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	msg := "Something went wrong"

	var httpErr interface{ HTTPStatusCode() int }
	if errors.As(err, &httpErr) {
		code = httpErr.HTTPStatusCode()
		msg = err.Error()
	}

	WriteJSON(w, code, HTTPError{
		Code:    code,
		Message: msg,
	})
}
