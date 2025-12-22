package http

import (
	"encoding/json"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service  service.UserService
	validate *validator.Validate
}

func NewHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		service:  svc,
		validate: validation.Get(),
	}
}

func (h *UserHandler) UsersGet(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Kullanıcılar getirilemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, users, http.StatusOK, "Kullanıcılar başarıyla getirildi")
}

func (h *UserHandler) UserPost(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := utils.ErrorResponse("VALIDATION_ERROR", "Geçersiz veri formatı", err.Error())
		utils.Return(w, http.StatusBadRequest, resp)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		resp := utils.ErrorResponse("VALIDATION_ERROR", "Geçersiz veri formatı", validation.FormatErr(err))
		utils.Return(w, http.StatusBadRequest, resp)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		resp := utils.ErrorResponse("DATABASE_ERROR", "Kullanıcı oluşturulamadı (İsim kullanımda olabilir)", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, user, http.StatusCreated, "Kullanıcı başarıyla eklendi")
}

func (h *UserHandler) UserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Kullanıcı silinemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	resp := utils.SuccessResponse(nil, "Kullanıcı başarıyla silindi")
	utils.Return(w, http.StatusOK, resp)
}
