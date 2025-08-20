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

type AuthService struct {
	userRepository *store.UserRepository
	jwtConfig      *security.JWTConfig
}

func NewAuthService(userRepo *store.UserRepository, jwtConfig *security.JWTConfig) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		jwtConfig:      jwtConfig,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	hashedPassword, err := security.GenerateHashFromPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword
	created, err := s.userRepository.Create(ctx, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				return nil, NewServiceError(http.StatusConflict, fmt.Sprintf("User with this email= [%s] already exists", user.Email), err)
			} else {
				return nil, NewServiceError(http.StatusConflict, fmt.Sprintf("User with this username= [%s] already exists", user.Username), err)
			}
		}

		return nil, NewServiceError(http.StatusInternalServerError, "Failed to create user", err)
	}
	return created, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.verifyCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	token, err := security.GenerateJWT(user.Id, s.jwtConfig)
	if err != nil {
		return nil, "", NewServiceError(http.StatusInternalServerError, "Generating jwt token failed", err)
	}

	return user, token, nil
}

func (s *AuthService) verifyCredentials(cxt context.Context, email, password string) (*model.User, error) {
	user, exists, err := s.userRepository.FindByEmail(cxt, email)

	if err != nil || !exists {
		return nil, NewServiceError(http.StatusBadRequest, "Invalid username or password", err)
	}
	if !security.CompareHashAndPassword(password, user.Password) {
		return nil, NewServiceError(http.StatusInternalServerError, "Password hashing failed", err)
	}
	return user, nil
}
