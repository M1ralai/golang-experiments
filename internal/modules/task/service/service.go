package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/events"
	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/outbox"
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
	userProvider domain.UserProvider
	outboxRepo   outbox.Repository
	logger       logger.Logger
}

func NewTaskService(
	taskRepo domain.TaskRepository,
	assignRepo domain.AssignmentRepository,
	activityRepo domain.ActivityRepository,
	userProvider domain.UserProvider,
	outboxRepo outbox.Repository,
	logger logger.Logger,
) TaskService {
	return &taskService{
		taskRepo:     taskRepo,
		assignRepo:   assignRepo,
		activityRepo: activityRepo,
		userProvider: userProvider,
		outboxRepo:   outboxRepo,
		logger:       logger,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req *domain.CreateTaskRequest) (*domain.Task, error) {
	userIDStr := utils.GetUserIDFromContext(ctx)
	createdBy, _ := uuid.Parse(userIDStr)
	if createdBy == uuid.Nil {
		createdBy = uuid.New()
	}

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

	userIDStr := utils.GetUserIDFromContext(ctx)
	userID, _ := uuid.Parse(userIDStr)
	if userID == uuid.Nil {
		userID = uuid.New()
	}
	activity := &domain.Activity{
		ID:        uuid.New(),
		TaskID:    uuid.MustParse(taskID),
		UserID:    userID,
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
	tx, err := s.assignRepo.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction", err, nil)
		return nil, err
	}
	defer tx.Rollback()

	assignment := &domain.TaskAssignment{
		ID:        uuid.New(),
		TaskID:    uuid.MustParse(taskID),
		UserID:    uuid.MustParse(req.UserID),
		CreatedAt: time.Now(),
	}

	if err := s.assignRepo.Create(ctx, tx, assignment); err != nil {
		s.logger.Error("Failed to assign task", err, map[string]interface{}{
			"task_id": taskID,
			"user_id": req.UserID,
		})
		return nil, err
	}

	activity := &domain.Activity{
		ID:        uuid.New(),
		TaskID:    assignment.TaskID,
		UserID:    assignment.UserID,
		Action:    domain.ActivityAssignmentAdded,
		CreatedAt: time.Now(),
	}
	_ = s.activityRepo.Create(ctx, activity)

	userInfo, err := s.userProvider.GetUserByID(assignment.UserID)
	if err != nil {
		s.logger.Error("Failed to get user info for event", err, map[string]interface{}{
			"user_id": req.UserID,
		})
		userInfo = &domain.UserInfo{
			ID:       assignment.UserID,
			Username: "Unknown",
			Email:    "",
		}
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to get task info for event", err, map[string]interface{}{
			"task_id": taskID,
		})
		task = &domain.Task{
			ID:    assignment.TaskID,
			Title: "Unknown Task",
		}
	}

	event := events.TaskAssignedEvent{
		TaskID:    taskID,
		TaskTitle: task.Title,
		UserID:    req.UserID,
		UserEmail: userInfo.Email,
		UserName:  userInfo.Username,
	}

	eventPayload, err := json.Marshal(event)
	if err != nil {
		s.logger.Error("Failed to marshal event", err, nil)
		return nil, err
	}

	outboxEvent := &outbox.OutboxEvent{
		AggregateType: "task",
		AggregateID:   assignment.TaskID,
		EventType:     events.TopicTaskAssigned,
		Payload:       eventPayload,
	}

	if err := s.outboxRepo.Create(ctx, tx, outboxEvent); err != nil {
		s.logger.Error("Failed to create outbox event", err, map[string]interface{}{
			"task_id": taskID,
			"user_id": req.UserID,
		})
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", err, nil)
		return nil, err
	}

	s.logger.Info("Task assigned", map[string]interface{}{
		"action":     "TASK_ASSIGN",
		"task_id":    taskID,
		"user_id":    req.UserID,
		"user_email": userInfo.Email,
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
