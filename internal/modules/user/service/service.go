package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/common/utils"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUsers() ([]domain.User, error)
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo   domain.UserRepository
	logger *logger.ZapLogger
}

func NewService(repo domain.UserRepository, logger *logger.ZapLogger) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) ListUsers() ([]domain.User, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("Failed to list users", err, nil)
		return nil, err
	}

	s.logger.Info("Users listed", map[string]interface{}{
		"action": "USER_LIST",
		"count":  len(users),
	})

	return users, nil
}

func (s *userService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		s.logger.Error("Failed to hash password", err, nil)
		return nil, err
	}

	user := &domain.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
		Ad:       req.Ad,
		Soyad:    req.Soyad,
		Telefon:  req.Telefon,
		Email:    req.Email,
	}

	err = s.repo.Create(user)
	if err != nil {
		s.logger.Error("Failed to create user", err, map[string]interface{}{
			"username": req.Username,
		})
		return nil, err
	}

	actorUsername := utils.GetUsernameFromContext(ctx)
	s.logger.Info("User created", map[string]interface{}{
		"action":       "USER_CREATE",
		"actor":        actorUsername,
		"new_username": user.Username,
		"role":         user.Role,
	})

	user.Password = ""
	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete user", err, map[string]interface{}{"user_id": id})
		return err
	}

	actorUsername := utils.GetUsernameFromContext(ctx)
	s.logger.Info("User deleted", map[string]interface{}{
		"action":  "USER_DELETE",
		"actor":   actorUsername,
		"user_id": id,
	})

	return nil
}
