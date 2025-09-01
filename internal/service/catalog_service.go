package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/errors"
	"arabic/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ICatalogService interface {
	GetAll(cxt context.Context) ([]*dto.CatalogResponse, error)
	Create(cxt context.Context, req *dto.CatalogCreateRequest) (int64, error)
	Delete(cxt context.Context, id int64) error
	Update(cxt context.Context, req *dto.CatalogUpdateRequest) error
	GetById(ctx context.Context, id int64) (*dto.GetCatalogByIdResponse, error)
}

type CatalogService struct {
	CatalogRepository repository.ICatalogRepository
}

func NewCatalogService(repo repository.ICatalogRepository) *CatalogService {
	return &CatalogService{CatalogRepository: repo}
}

func (c *CatalogService) Delete(ctx context.Context, id int64) error {
	ok, err := c.CatalogRepository.Delete(ctx, id)

	if err != nil {
		logger.Log.Error("CatalogService -> Delete -> err -> " + err.Error())
		return errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	if !ok {
		return errors.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Cant delete catalog item by id: %d, please check provided id", id), nil)
	}

	return nil
}

func (c *CatalogService) Create(cxt context.Context, req *dto.CatalogCreateRequest) (int64, error) {
	item, err := c.CatalogRepository.Create(cxt, &model.Catalog{
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
		logger.Log.Error("CatalogService -> Create -> err -> " + err.Error())
		return 0, errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	return item.Id, nil
}

func (c *CatalogService) Update(cxt context.Context, req *dto.CatalogUpdateRequest) error {
	query, values := c.prepareQueryForUpdate(req)
	ok, err := c.CatalogRepository.Update(cxt, query, values)

	if err != nil {
		logger.Log.Error("CatalogService -> Update -> err -> " + err.Error())
		return errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	if !ok {
		logger.Log.Error(fmt.Sprintf("CatalogService -> Update -> err -> "+"Cant update catalog item id: %d", req.Id))
		return errors.NewServiceError(http.StatusBadRequest, errors.ErrorNotFoundById, nil)
	}

	return nil
}

func (c *CatalogService) GetById(ctx context.Context, id int64) (*dto.GetCatalogByIdResponse, error) {
	item, ok, err := c.CatalogRepository.FindById(ctx, id)

	if err != nil {
		logger.Log.Error("CatalogService -> GetById -> err -> " + err.Error())
		return nil, errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	if !ok {
		return nil, errors.NewServiceError(http.StatusBadRequest, errors.ErrorNotFoundById, nil)
	}

	resp := &dto.GetCatalogByIdResponse{
		Id:          item.Id,
		Description: item.Description,
		Sku:         item.Sku,
		CatalogBase: dto.CatalogBase{
			Name:            item.Name,
			Price:           item.Price,
			Amount:          item.Amount,
			CategoryId:      item.CategoryId,
			DiscountPercent: item.DiscountPercent,
		},
	}

	return resp, nil
}

func (c *CatalogService) GetAll(cxt context.Context) ([]*dto.CatalogResponse, error) {

	catalogItems, err := c.CatalogRepository.FindAll(cxt)

	if err != nil {
		logger.Log.Error("CatalogService -> GetAll -> err -> " + err.Error())
		return nil, errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil)
	}

	var catalogResp []*dto.CatalogResponse

	for _, item := range catalogItems {
		catalogResp = append(catalogResp, &dto.CatalogResponse{
			Id: item.Id,
			CatalogBase: dto.CatalogBase{
				Name:            item.Name,
				Price:           item.Price,
				CategoryId:      item.CategoryId,
				Amount:          item.Amount,
				DiscountPercent: item.DiscountPercent,
			},
		})
	}

	return catalogResp, nil
}

// Покрыть тестом
func (c *CatalogService) prepareQueryForUpdate(req *dto.CatalogUpdateRequest) (string, []any) {
	fieldsForUpdate := make([]string, 0)
	values := make([]any, 0)

	query := ""
	paramIndex := 1

	values = append(values, req.Id)
	paramIndex++

	if req.Name != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("name = $%d", paramIndex))
		values = append(values, *req.Name)
		paramIndex++
	}

	if req.Description != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("description = $%d", paramIndex))
		values = append(values, *req.Description)
		paramIndex++
	}

	if req.Price != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("price = $%d", paramIndex))
		values = append(values, *req.Price)
		paramIndex++
	}

	if req.DiscountPercent != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("discount_percent = $%d", paramIndex))
		values = append(values, *req.DiscountPercent)
		paramIndex++
	}

	if req.Amount != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("amount = $%d", paramIndex))
		values = append(values, *req.Amount)
		paramIndex++
	}

	if req.CategoryId != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("category_id = $%d", paramIndex))
		values = append(values, *req.CategoryId)
		paramIndex++
	}

	if req.Sku != nil {
		fieldsForUpdate = append(fieldsForUpdate, fmt.Sprintf("sku = $%d", paramIndex))
		values = append(values, *req.Sku)
		paramIndex++
	}

	fieldsLength := len(fieldsForUpdate)
	for idx, val := range fieldsForUpdate {
		query += val
		if idx+1 != fieldsLength {
			query += ", "
		}
	}

	return query, values
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
