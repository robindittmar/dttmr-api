package handler

import (
	"log/slog"
	"net/http"

	"github.com/robindittmar/dttmr-api/internal/api/request"
	"github.com/robindittmar/dttmr-api/internal/api/response"
	"github.com/robindittmar/dttmr-api/internal/domain"
)

type AuthHandler struct {
	AuthService *domain.AuthService
}

func NewAuthHandler(authService *domain.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := request.DecodeLogin(r)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode login payload", slog.Any("error", err))
		response.Error(ctx, w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	token, err := h.AuthService.Login(ctx, payload.Email, payload.Password)
	if err != nil {
		slog.ErrorContext(ctx, "failed to login", slog.Any("error", err))
		response.Error(ctx, w, http.StatusInternalServerError, "failed to login")
		return
	}

	response.JSON(ctx, w, http.StatusOK, token)
}
