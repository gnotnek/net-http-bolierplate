package resp

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
)

type Meta struct {
	Page      int `json:"page"`
	PageTotal int `json:"page_total"`
	Total     int `json:"total"`
}

type DataPaginate struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type HTTPError struct {
	StatusCode int   `json:"code"`
	Message    error `json:"message"`
}

type Empty struct{}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonData, _ := json.Marshal(data)
	w.Write(jsonData)
}

func WriteJSONWithPaginateResponse(w http.ResponseWriter, statusCode int, data interface{}, page, total, limit int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	meta := Meta{
		Page:      page,
		PageTotal: totalPage,
		Total:     total,
	}

	jsonData, _ := json.Marshal(DataPaginate{
		Data: data,
		Meta: meta,
	})

	w.Write(jsonData)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	code := http.StatusInternalServerError
	msg := "Something went wrong"

	var httpErr interface{ HTTPStatusCode() int }
	if errors.As(err, &httpErr) {
		code = httpErr.HTTPStatusCode()
		msg = err.Error()
	}

	errResponse := HTTPError{
		StatusCode: code,
		Message:    errors.New(msg),
	}

	response, _ := json.Marshal(errResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
