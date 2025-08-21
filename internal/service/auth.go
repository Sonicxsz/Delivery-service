package service

import (
	"arabic/internal/model"
	security "arabic/internal/security/auth"
	. "arabic/pkg/errors"
	"arabic/store"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type UserService struct {
	userRepository *store.UserRepository
	jwtConfig      *security.JWTConfig
}

type AuthService interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, string, error)
}

func NewAuthService(userRepo *store.UserRepository, jwtConfig *security.JWTConfig) *UserService {
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
		return s.handleDuplicateErrorMessage(err, user)
	}

	return nil, NewServiceError(http.StatusInternalServerError, "Failed to create user", err)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.verifyCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := security.GenerateJWT(user.Id, s.jwtConfig)
	if err != nil {
		return nil, "", NewServiceError(http.StatusInternalServerError, "Something went wrong. pls try later", err)
	}

	return user, token, nil
}

func (s *UserService) verifyCredentials(cxt context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepository.FindByEmail(cxt, email)

	if err != nil || !security.CompareHashAndPassword(password, user.Password) {
		return nil, NewServiceError(http.StatusBadRequest, "Invalid username or password", err)
	}

	return user, nil
}

func (s *UserService) handleDuplicateErrorMessage(err error, user *model.User) (*model.User, error) {
	if strings.Contains(err.Error(), "email") {
		return nil, NewServiceError(http.StatusConflict, fmt.Sprintf("User with this email= [%s] already exists", user.Email), err)
	}
	return nil, NewServiceError(http.StatusConflict, fmt.Sprintf("User with this username= [%s] already exists", user.Username), err)
}
