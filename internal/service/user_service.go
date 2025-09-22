package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/customError"
	"arabic/pkg/logger"
	"arabic/pkg/security/auth"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, string, error)
	GetUser(ctx context.Context, email string) (*dto.UserGetResponse, error)
}

type UserService struct {
	userRepository repository.IUserRepository
	jwtConfig      *security.JWTConfig
}

func NewUserService(userRepo repository.IUserRepository, jwtConfig *security.JWTConfig) *UserService {
	return &UserService{
		userRepository: userRepo,
		jwtConfig:      jwtConfig,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	hashedPassword, err := security.GenerateHashFromPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	created, err := s.userRepository.Create(ctx, user)

	if err == nil {
		return created, nil
	}

	if strings.Contains(err.Error(), "duplicate key") {
		return nil, s.handleDuplicateErrorMessage(err, user)
	}

	return nil, customError.NewServiceError(http.StatusInternalServerError, "Failed to create user", err)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.verifyCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := security.GenerateJWT(user.Email, s.jwtConfig)
	if err != nil {
		return nil, "", customError.NewServiceError(http.StatusInternalServerError, "Something went wrong. pls try later", err)
	}

	return user, token, nil
}

func (s *UserService) GetUser(ctx context.Context, email string) (*dto.UserGetResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)

	if err != nil {
		logger.Log.Error("UserService -> GetUser -> err -> " + err.Error())
		return nil, customError.NewServiceError(http.StatusBadRequest, "Cant find user by provided email", err)
	}

	resp := dto.UserGetResponse{
		Email:    user.Email,
		Username: user.Username,
	}

	return &resp, nil
}

func (s *UserService) verifyCredentials(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)

	if err != nil || !security.CompareHashAndPassword(password, user.Password) {
		return nil, customError.NewServiceError(http.StatusBadRequest, "Invalid username or password", err)
	}

	return user, nil
}

func (s *UserService) handleDuplicateErrorMessage(err error, user *model.User) error {
	if strings.Contains(err.Error(), "email") {
		return customError.NewServiceError(http.StatusConflict, fmt.Sprintf("User with this email= [%s] already exists", user.Email), err)
	}
	return customError.NewServiceError(http.StatusConflict, fmt.Sprintf("User with this username= [%s] already exists", user.Username), err)
}
