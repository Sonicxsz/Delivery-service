package service

import (
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/customError"
	"arabic/pkg/security/auth"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, string, error)
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

	token, err := security.GenerateJWT(strconv.FormatInt(user.Id, 10), s.jwtConfig)
	if err != nil {
		return nil, "", customError.NewServiceError(http.StatusInternalServerError, "Something went wrong. pls try later", err)
	}

	return user, token, nil
}

func (s *UserService) verifyCredentials(cxt context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepository.FindByEmail(cxt, email)

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
