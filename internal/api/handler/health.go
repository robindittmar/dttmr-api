package handler

import (
	"net/http"

	"github.com/robindittmar/dttmr-api/internal/api/response"
)

type healthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{
		Status: "ok",
	}
	response.JSON(w, http.StatusOK, resp)
}
