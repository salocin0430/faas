package service

import (
	"context"
	"faas/internal/features/users/application/dto"
	"faas/internal/features/users/domain/entity"
	"faas/internal/features/users/domain/repository"
	"faas/internal/shared/infrastructure/config"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{repo: repo, cfg: cfg}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	log.Printf("Password before hash: %s", req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}
	log.Printf("Generated hash: %s", string(hashedPassword))

	user := &entity.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Password:  string(hashedPassword),
		Role:      req.Role,
		CreatedAt: time.Now(),
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Verification failed immediately after hashing: %v", err)
		return nil, fmt.Errorf("password hashing verification failed")
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, err
	}

	// Generar JWT
	token, err := s.generateJWT(user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return nil, err
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (s *UserService) generateJWT(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"user": user.Username,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"key":  s.cfg.ConsumerKey,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *UserService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		response[i] = &dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}
	}

	return response, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
