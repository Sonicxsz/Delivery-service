package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/errors"
	"context"
	"fmt"
	"net/http"
)

type ICategoryService interface {
	GetAll(cxt context.Context) ([]*dto.CategoryResponse, error)
	Create(cxt context.Context, req *dto.CategoryRequest) (*dto.CategoryResponse, error)
	Delete(cxt context.Context, id int64) error
}

type CategoryService struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryService(categoryRepository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *CategoryService) GetAll(cxt context.Context) ([]*dto.CategoryResponse, error) {
	all, err := s.categoryRepository.FindAll(cxt)
	if err != nil {
		return nil, err
	}

	var response []*dto.CategoryResponse
	for _, category := range all {
		response = append(response, &dto.CategoryResponse{
			Id:   category.Id,
			Name: category.Name,
			Code: category.Code,
		})
	}

	return response, nil
}

func (s *CategoryService) Create(cxt context.Context, req *dto.CategoryRequest) (*dto.CategoryResponse, error) {
	category := &model.Category{
		Name: req.Name,
		Code: req.Code,
	}

	created, err := s.categoryRepository.Create(cxt, category)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		Id:   created.Id,
		Name: created.Name,
		Code: created.Code,
	}, nil
}

func (s *CategoryService) Delete(cxt context.Context, id int64) error {
	category, err := s.categoryRepository.Delete(cxt, id)
	if err != nil {
		return err
	} else if category.RowsAffected() == 0 {
		return errors.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Entity not found with id %d", id), err)
	}
	return nil
}
