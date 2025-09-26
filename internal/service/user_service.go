package service

import (
	"arabic/internal/dto"
	"arabic/internal/repository"
	"arabic/pkg/customError"
	"arabic/pkg/logger"
	"arabic/pkg/queryBuilder"
	"arabic/pkg/security/auth"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type IUserService interface {
	CreateUser(ctx context.Context, user *dto.UserCreateRequest) error
	Login(ctx context.Context, email, password string) (*dto.UserGetResponse, string, error)
	GetUser(ctx context.Context, email string) (*dto.UserGetResponse, error)
	UpdateUserInfo(ctx context.Context, req *dto.UserUpdateRequest) error
	UpdateUserAddress(cxt context.Context, req *dto.UserAddressUpdateRequest) error
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

func (s *UserService) CreateUser(ctx context.Context, user *dto.UserCreateRequest) error {
	hashedPassword, err := security.GenerateHashFromPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	err = s.userRepository.Create(ctx, user)

	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key") {
		return s.handleDuplicateErrorMessage(err, user)
	}

	return customError.NewServiceError(http.StatusInternalServerError, "Failed to create user", err)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*dto.UserGetResponse, string, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)

	if err != nil {
		security.CompareHashAndPassword("Dummy-password-for-time", password)
		return nil, "", customError.NewServiceError(http.StatusBadRequest, "Invalid username or password", err)
	}

	ok := security.CompareHashAndPassword(password, user.Password)
	if !ok {
		return nil, "", customError.NewServiceError(http.StatusBadRequest, "Invalid username or password", err)
	}

	token, err := security.GenerateJWT(user.Email, user.Id, s.jwtConfig)
	if err != nil {
		return nil, "", customError.NewServiceError(http.StatusInternalServerError, "Something went wrong. pls try later", err)
	}

	resp := &dto.UserGetResponse{
		Email:      user.Email,
		SecondName: user.SecondName,
		FirstName:  user.FirstName,
		Id:         user.Id,
		Username:   user.Username,
	}

	return resp, token, nil
}

func (s *UserService) GetUser(ctx context.Context, email string) (*dto.UserGetResponse, error) {
	user, err := s.userRepository.FindByEmail(ctx, email)

	if err != nil {
		logger.Log.Error("UserService -> GetUser -> err -> " + err.Error())
		return nil, customError.NewServiceError(http.StatusBadRequest, "Cant find user by provided email", err)
	}

	resp := dto.UserGetResponse{
		Email:      user.Email,
		Username:   user.Username,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Id:         user.Id,
	}

	return &resp, nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *dto.UserUpdateRequest) error {
	qb := queryBuilder.NewQueryBuilder(true).
		Set("first_name", req.FirstName).
		Set("second_name", req.SecondName)

	query, values := qb.BuildUpdateQuery("public.users", "id", req.Id)

	// Юзер не предоставил данные для измения(пустой запрос)
	if len(values) < 2 {
		return customError.NewServiceError(http.StatusBadRequest, customError.ErrorWrongPayload, nil)
	}

	ok, err := s.userRepository.Update(ctx, query, values)

	if err != nil {
		logger.Log.Error("UserService -> UpdateUser -> err -> " + err.Error())
		return customError.NewServiceError(http.StatusInternalServerError, customError.Error500, err)
	}

	if !ok {
		return customError.NewServiceError(http.StatusInternalServerError, customError.ErrorNotFoundById, err)
	}

	return nil
}

func (s *UserService) UpdateUserAddress(cxt context.Context, req *dto.UserAddressUpdateRequest) error {
	// Подумать нужно ли перенести формирование sql запроса на уровень репозитория
	qb := queryBuilder.NewQueryBuilder(false).
		Set("region", req.Region).
		Set("city", req.City).
		Set("street", req.Street).
		Set("house", req.House).
		Set("apartment", req.Apartment)

	query, values := qb.BuildUpdateQuery("public.users", "id", req.Id)

	if values == nil {
		return customError.NewServiceError(http.StatusInternalServerError, customError.ErrorWrongPayload, nil)
	}

	// Нужно будет проверить корректность предоставленного адреса
	ok, err := s.userRepository.Update(cxt, query, values)

	if err != nil {
		logger.Log.Error("UserService -> UpdateUserAddress -> err -> " + err.Error())
		return customError.NewServiceError(http.StatusInternalServerError, customError.Error500, err)
	}

	if !ok {
		return customError.NewServiceError(http.StatusInternalServerError, customError.ErrorNotFoundById, err)
	}

	return nil
}

func (s *UserService) handleDuplicateErrorMessage(err error, user *dto.UserCreateRequest) error {
	if strings.Contains(err.Error(), "email") {
		return customError.NewServiceError(http.StatusConflict, fmt.Sprintf("User with this email= [%s] already exists", user.Email), err)
	}
	return customError.NewServiceError(http.StatusConflict, fmt.Sprintf("User with this username= [%s] already exists", user.Username), err)
}
