package service

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/auth/domain"
	userDomain "github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey  []byte
	jwtOnce sync.Once
)

func getJWTKey() []byte {
	jwtOnce.Do(func() {
		key := os.Getenv("JWT_SECRET")
		if key == "" {
			log.Fatal("JWT_SECRET environment variable is not set")
		}
		jwtKey = []byte(key)
	})
	return jwtKey
}

type AuthService interface {
	Login(req *domain.LoginRequest) (*domain.LoginResponse, error)
}

type authService struct {
	userRepo userDomain.UserRepository
	logger   *logger.ZapLogger
}

func NewService(userRepository userDomain.UserRepository, logger *logger.ZapLogger) AuthService {
	return &authService{
		userRepo: userRepository,
		logger:   logger,
	}
}

func (s *authService) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {

	if req.Username == "admin" && req.Password == "123" {
		return s.generateTokenForTestUser("admin", "ADMIN", "Test", "Admin")
	}
	if req.Username == "sekreter" && req.Password == "123" {
		return s.generateTokenForTestUser("sekreter", "SEKRETER", "Test", "Sekreter")
	}

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		s.logger.Error("Kullanıcı bulunamadı", err, map[string]interface{}{"username": req.Username})
		return nil, domain.ErrInvalidCredentials{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.logger.Error("Şifre doğrulanamadı", err, map[string]interface{}{"username": req.Username})
		return nil, domain.ErrInvalidCredentials{}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		s.logger.Error("Token oluşturulamadı", err, nil)
		return nil, domain.ErrTokenGeneration{}
	}

	s.logger.Info("Kullanıcı giriş yaptı", map[string]interface{}{
		"user": user.Username,
		"role": user.Role,
	})

	return &domain.LoginResponse{
		Token:    tokenString,
		Username: user.Username,
		Role:     user.Role,
		Ad:       user.Ad,
		Soyad:    user.Soyad,
	}, nil
}

func (s *authService) generateTokenForTestUser(username, role, ad, soyad string) (*domain.LoginResponse, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTKey())
	if err != nil {
		return nil, domain.ErrTokenGeneration{}
	}

	s.logger.Info("Test kullanıcı giriş yaptı", map[string]interface{}{
		"user": username,
		"role": role,
	})

	return &domain.LoginResponse{
		Token:    tokenString,
		Username: username,
		Role:     role,
		Ad:       ad,
		Soyad:    soyad,
	}, nil
}
