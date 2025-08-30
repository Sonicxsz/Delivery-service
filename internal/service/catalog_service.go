package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/errors"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ICatalogService interface {
	GetAll(cxt context.Context) ([]*dto.CatalogResponse, error)
	Create(cxt context.Context, req *dto.CatalogCreateRequest) (int64, error)
	Delete(cxt context.Context, id int64) error
}

type CatalogService struct {
	catalogRepository repository.ICatalogRepository
}

func NewCatalogService(repo repository.ICatalogRepository) *CatalogService {
	return &CatalogService{catalogRepository: repo}
}

func (c *CatalogService) Delete(ctx context.Context, id int64) error {
	ok, err := c.catalogRepository.Delete(ctx, id)

	if err != nil {
		return errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	if !ok {
		return errors.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Cant delete catalog item by id: %d, please check provided id", id), nil)
	}

	return nil
}

func (c *CatalogService) Create(cxt context.Context, req *dto.CatalogCreateRequest) (int64, error) {
	item, err := c.catalogRepository.Create(cxt, &model.Catalog{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Amount:          req.Amount,
		Sku:             req.Sku,
		DiscountPercent: req.DiscountPercent,
		CategoryId:      req.CategoryId,
	})

	if err != nil && c.hasDuplicates(err) {
		return 0, c.getCatalogUniqFieldError(err, item)
	}

	if err != nil {
		return 0, errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	return item.Id, nil
}

func (c *CatalogService) GetAll(cxt context.Context) ([]*dto.CatalogResponse, error) {

	catalogItems, err := c.catalogRepository.FindAll(cxt)

	if err != nil {
		return nil, err
	}

	var catalogResp []*dto.CatalogResponse

	for _, item := range catalogItems {
		catalogResp = append(catalogResp, &dto.CatalogResponse{
			Id:              item.Id,
			Name:            item.Name,
			Price:           item.Price,
			CategoryId:      item.CategoryId,
			Amount:          item.Amount,
			DiscountPercent: item.DiscountPercent,
		})
	}

	return catalogResp, nil
}

func (c *CatalogService) hasDuplicates(err error) bool {
	return strings.Contains(err.Error(), "duplicate")
}

func (c *CatalogService) getCatalogUniqFieldError(err error, catalog *model.Catalog) error {
	if strings.Contains(err.Error(), "sku") {
		return errors.NewServiceError(http.StatusConflict, fmt.Sprintf("Provided sku: %s already exist", catalog.Sku), err)
	}

	return errors.NewServiceError(http.StatusConflict, "Please check provided data", err)
}
