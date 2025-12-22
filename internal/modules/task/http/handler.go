package http

import (
	"encoding/json"
	"net/http"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/validation"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/service"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service  service.TaskService
	validate *validator.Validate
}

func NewHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{
		service:  svc,
		validate: validation.Get(),
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTaskRequest
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

	task, err := h.service.CreateTask(r.Context(), &req)
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task oluşturulamadı", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, task, http.StatusCreated, "Task başarıyla oluşturuldu")
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.ListTasks(r.Context())
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task listesi getirilemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, tasks, http.StatusOK, "Task listesi başarıyla getirildi")
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	task, err := h.service.GetTask(r.Context(), taskID)
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task getirilemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	if task == nil {
		resp := utils.ErrorResponse("NOT_FOUND", "Task bulunamadı", "")
		utils.Return(w, http.StatusNotFound, resp)
		return
	}

	utils.WriteJson(w, task, http.StatusOK, "Task başarıyla getirildi")
}

func (h *TaskHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var req domain.UpdateStatusRequest
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

	if err := h.service.UpdateTaskStatus(r.Context(), taskID, &req); err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task durumu güncellenemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	resp := utils.SuccessResponse(nil, "Task durumu başarıyla güncellendi")
	utils.Return(w, http.StatusOK, resp)
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	var req domain.AssignTaskRequest
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

	assignment, err := h.service.AssignTask(r.Context(), taskID, &req)
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task ataması yapılamadı", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, assignment, http.StatusCreated, "Task başarıyla atandı")
}

func (h *TaskHandler) UnassignTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assignmentID := vars["id"]

	if err := h.service.UnassignTask(r.Context(), assignmentID); err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task ataması kaldırılamadı", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	resp := utils.SuccessResponse(nil, "Task ataması başarıyla kaldırıldı")
	utils.Return(w, http.StatusOK, resp)
}

func (h *TaskHandler) GetTaskAssignments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	assignments, err := h.service.GetTaskAssignments(r.Context(), taskID)
	if err != nil {
		resp := utils.ErrorResponse("INTERNAL_ERROR", "Task atamaları getirilemedi", err.Error())
		utils.Return(w, http.StatusInternalServerError, resp)
		return
	}

	utils.WriteJson(w, assignments, http.StatusOK, "Task atamaları başarıyla getirildi")
}
