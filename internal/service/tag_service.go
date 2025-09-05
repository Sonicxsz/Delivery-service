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

type TagService struct {
	TagRepository repository.ITagRepository
}

type ITagService interface {
	GetAll(cxt context.Context) ([]*dto.TagResponse, error)
	Create(cxt context.Context, req *dto.TagRequest) (*dto.TagResponse, error)
	Delete(cxt context.Context, id int64) error
}

func NewTagService(tagRepository *repository.TagRepository) *TagService {
	return &TagService{
		TagRepository: tagRepository,
	}
}

func (s *TagService) GetAll(cxt context.Context) ([]*dto.TagResponse, error) {
	all, err := s.TagRepository.FindAll(cxt)
	if err != nil {
		return nil, err
	}

	var response []*dto.TagResponse
	for _, tag := range all {
		response = append(response, &dto.TagResponse{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}

	return response, nil
}

func (s *TagService) Create(cxt context.Context, req *dto.TagRequest) (*dto.TagResponse, error) {
	tag := &model.Tag{
		Name: req.Name,
	}

	created, err := s.TagRepository.Create(cxt, tag)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		Id:   created.Id,
		Name: created.Name,
	}, nil
}

func (s *TagService) Delete(cxt context.Context, id int64) error {
	tag, err := s.TagRepository.Delete(cxt, id)
	if err != nil {
		return err
	} else if tag.RowsAffected() == 0 {
		return errors.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Entity not found with id %d", id), err)
	}
	return nil
}
