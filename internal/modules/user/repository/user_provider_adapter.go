package repository

import (
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	userDomain "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/google/uuid"
)

type UserProviderAdapter struct {
	userRepo userDomain.UserRepository
}

func NewUserProviderAdapter(userRepo userDomain.UserRepository) domain.UserProvider {
	return &UserProviderAdapter{
		userRepo: userRepo,
	}
}

func (a *UserProviderAdapter) GetUserByID(userID uuid.UUID) (*domain.UserInfo, error) {
	user, err := a.userRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
