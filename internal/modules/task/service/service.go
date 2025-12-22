package service

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	"github.com/google/uuid"
)

type TaskService interface {
	CreateTask(ctx context.Context, req *domain.CreateTaskRequest) (*domain.Task, error)
	GetTask(ctx context.Context, taskID string) (*domain.Task, error)
	ListTasks(ctx context.Context) ([]domain.Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, req *domain.UpdateStatusRequest) error

	AssignTask(ctx context.Context, taskID string, req *domain.AssignTaskRequest) (*domain.TaskAssignment, error)
	UnassignTask(ctx context.Context, assignmentID string) error
	GetTaskAssignments(ctx context.Context, taskID string) ([]domain.TaskAssignment, error)
}

type taskService struct {
	taskRepo     domain.TaskRepository
	assignRepo   domain.AssignmentRepository
	activityRepo domain.ActivityRepository
	logger       *logger.ZapLogger
}

func NewTaskService(
	taskRepo domain.TaskRepository,
	assignRepo domain.AssignmentRepository,
	activityRepo domain.ActivityRepository,
	logger *logger.ZapLogger,
) TaskService {
	return &taskService{
		taskRepo:     taskRepo,
		assignRepo:   assignRepo,
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req *domain.CreateTaskRequest) (*domain.Task, error) {
	createdBy := uuid.New() // TODO: Get from context when auth is implemented

	task := &domain.Task{
		ID:        uuid.New(),
		Title:     req.Title,
		Status:    domain.TaskStatusTodo,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.taskRepo.Create(ctx, task)
	if err != nil {
		s.logger.Error("Failed to create task", err, map[string]interface{}{
			"title": req.Title,
		})
		return nil, err
	}

	// Log activity
	activity := &domain.Activity{
		ID:        uuid.New(),
		TaskID:    task.ID,
		UserID:    createdBy,
		Action:    domain.ActivityTaskCreated,
		CreatedAt: time.Now(),
	}
	_ = s.activityRepo.Create(ctx, activity)

	s.logger.Info("Task created", map[string]interface{}{
		"action":  "TASK_CREATE",
		"task_id": task.ID.String(),
		"title":   task.Title,
	})

	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to get task", err, map[string]interface{}{
			"task_id": taskID,
		})
		return nil, err
	}
	return task, nil
}

func (s *taskService) ListTasks(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.taskRepo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to list tasks", err, nil)
		return nil, err
	}

	s.logger.Info("Tasks listed", map[string]interface{}{
		"action": "TASK_LIST",
		"count":  len(tasks),
	})

	return tasks, nil
}

func (s *taskService) UpdateTaskStatus(ctx context.Context, taskID string, req *domain.UpdateStatusRequest) error {
	err := s.taskRepo.UpdateStatus(ctx, taskID, req.Status)
	if err != nil {
		s.logger.Error("Failed to update task status", err, map[string]interface{}{
			"task_id": taskID,
			"status":  req.Status,
		})
		return err
	}

	// Log activity
	activity := &domain.Activity{
		ID:        uuid.New(),
		TaskID:    uuid.MustParse(taskID),
		UserID:    uuid.New(), // TODO: Get from context
		Action:    domain.ActivityTaskStatusChanged,
		CreatedAt: time.Now(),
	}
	_ = s.activityRepo.Create(ctx, activity)

	s.logger.Info("Task status updated", map[string]interface{}{
		"action":     "TASK_STATUS_UPDATE",
		"task_id":    taskID,
		"new_status": req.Status,
	})

	return nil
}

func (s *taskService) AssignTask(ctx context.Context, taskID string, req *domain.AssignTaskRequest) (*domain.TaskAssignment, error) {
	assignment := &domain.TaskAssignment{
		ID:        uuid.New(),
		TaskID:    uuid.MustParse(taskID),
		UserID:    uuid.MustParse(req.UserID),
		CreatedAt: time.Now(),
	}

	err := s.assignRepo.Create(ctx, assignment)
	if err != nil {
		s.logger.Error("Failed to assign task", err, map[string]interface{}{
			"task_id": taskID,
			"user_id": req.UserID,
		})
		return nil, err
	}

	// Log activity
	activity := &domain.Activity{
		ID:        uuid.New(),
		TaskID:    assignment.TaskID,
		UserID:    assignment.UserID,
		Action:    domain.ActivityAssignmentAdded,
		CreatedAt: time.Now(),
	}
	_ = s.activityRepo.Create(ctx, activity)

	s.logger.Info("Task assigned", map[string]interface{}{
		"action":  "TASK_ASSIGN",
		"task_id": taskID,
		"user_id": req.UserID,
	})

	return assignment, nil
}

func (s *taskService) UnassignTask(ctx context.Context, assignmentID string) error {
	err := s.assignRepo.Delete(ctx, assignmentID)
	if err != nil {
		s.logger.Error("Failed to unassign task", err, map[string]interface{}{
			"assignment_id": assignmentID,
		})
		return err
	}

	s.logger.Info("Task unassigned", map[string]interface{}{
		"action":        "TASK_UNASSIGN",
		"assignment_id": assignmentID,
	})

	return nil
}

func (s *taskService) GetTaskAssignments(ctx context.Context, taskID string) ([]domain.TaskAssignment, error) {
	assignments, err := s.assignRepo.GetByTask(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to get task assignments", err, map[string]interface{}{
			"task_id": taskID,
		})
		return nil, err
	}
	return assignments, nil
}
