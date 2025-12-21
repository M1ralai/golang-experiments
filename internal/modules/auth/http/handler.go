package http

import (
	"encoding/json"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/service"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  service.AuthService
	validate *validator.Validate
}

func NewHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  svc,
		validate: validation.Get(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := utils.ErrorResponse("INVALID_INPUT", "Geçersiz veri formatı", err.Error())
		utils.Return(w, http.StatusBadRequest, resp)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		resp := utils.ErrorResponse("VALIDATION_ERROR", "Geçersiz veri formatı", validation.FormatErr(err))
		utils.Return(w, http.StatusBadRequest, resp)
		return
	}

	loginResp, err := h.service.Login(&req)
	if err != nil {
		if _, ok := err.(domain.ErrInvalidCredentials); ok {
			resp := utils.ErrorResponse("UNAUTHORIZED", "Hatalı kullanıcı adı veya şifre", err.Error())
			utils.Return(w, http.StatusUnauthorized, resp)
			return
		}
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Sunucu hatası", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, loginResp, http.StatusOK, "Giriş başarılı")
}
