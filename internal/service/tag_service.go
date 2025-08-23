package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/store"
	"context"
)

type TagService struct {
	tagRepository store.TagRepository
}

type TagServiceI interface {
	FindAll(ctx context.Context) ([]model.Tag, error)
}

func NewTagService(tagRepository *store.TagRepository) *TagService {
	return &TagService{
		tagRepository: *tagRepository,
	}
}

func (s *TagService) FindAll(cxt context.Context) ([]model.Tag, error) {
	return s.tagRepository.FindAll(cxt)
}

func (s *TagService) CreateTag(cxt context.Context, req *dto.TagRequest) (*dto.TagRequest, error) {
	return s.tagRepository.Create(cxt, req)
}
