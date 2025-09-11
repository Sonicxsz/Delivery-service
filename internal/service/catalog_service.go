package service

import (
	"arabic/internal/dto"
	"arabic/internal/model"
	"arabic/internal/repository"
	"arabic/pkg/customError"
	"arabic/pkg/fs"
	"arabic/pkg/logger"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type ICatalogService interface {
	GetAll(cxt context.Context, imagePrefix string) ([]*dto.CatalogResponse, error)
	Create(cxt context.Context, req *dto.CatalogCreateRequest) (uint, error)
	Delete(cxt context.Context, id uint) error
	Update(cxt context.Context, req *dto.CatalogUpdateRequest) error
	GetById(ctx context.Context, id uint, imagePrefix string) (*dto.CatalogResponse, error)
	AddImage(cxt context.Context, req *dto.AddImageRequest, fs fs.IFileSystemImage) (string, error)
}

type CatalogService struct {
	CatalogRepository repository.ICatalogRepository
}

func NewCatalogService(repo repository.ICatalogRepository) *CatalogService {
	return &CatalogService{CatalogRepository: repo}
}

func (c *CatalogService) Delete(ctx context.Context, id uint) error {
	ok, err := c.CatalogRepository.Delete(ctx, id)

	if err != nil {
		logger.Log.Error("CatalogService -> Delete -> err -> " + err.Error())
		return customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	if !ok {
		return customError.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Cant delete catalog item by id: %d, please check provided id", id), nil)
	}

	return nil
}

func (c *CatalogService) Create(cxt context.Context, req *dto.CatalogCreateRequest) (uint, error) {
	item, err := c.CatalogRepository.Create(cxt, &model.Catalog{
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Amount:          req.Amount,
		Sku:             req.Sku,
		DiscountPercent: req.DiscountPercent,
		CategoryId:      req.CategoryId,
	})

	if err != nil && isDuplicateError(err) {
		return 0, c.getCatalogUniqFieldError(err, item)
	}

	if err != nil {
		logger.Log.Error("CatalogService -> Create -> err -> " + err.Error())
		return 0, customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	return item.Id, nil
}

func (c *CatalogService) Update(cxt context.Context, req *dto.CatalogUpdateRequest) error {
	query, values := c.prepareQueryForUpdate(req)
	ok, err := c.CatalogRepository.Update(cxt, query, values)

	if err != nil {
		logger.Log.Error("CatalogService -> Update -> err -> " + err.Error())
		return customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	if !ok {
		logger.Log.Error(fmt.Sprintf("CatalogService -> Update -> err -> "+"Cant update catalog item id: %d", req.Id))
		return customError.NewServiceError(http.StatusBadRequest, customError.ErrorNotFoundById, nil)
	}

	return nil
}

func (c *CatalogService) GetById(ctx context.Context, id uint, imagePrefix string) (*dto.CatalogResponse, error) {
	item, ok, err := c.CatalogRepository.FindById(ctx, id)

	if err != nil {
		logger.Log.Error("CatalogService -> GetById -> err -> " + err.Error())
		return nil, customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	if !ok {
		return nil, customError.NewServiceError(http.StatusBadRequest, customError.ErrorNotFoundById, nil)
	}

	resp := item.ToResponse(imagePrefix)

	return resp, nil
}

func (c *CatalogService) AddImage(cxt context.Context, req *dto.AddImageRequest, fs fs.IFileSystemImage) (string, error) {

	extension, err := fs.GetImageExtension(&req.Image)

	if err != nil {
		logger.Log.Error(err.Error())
		return "", customError.NewServiceError(http.StatusBadRequest, "Image extension not found. Provide correct data", nil)
	}

	ok := fs.IsSupportingExtension(extension)

	if !ok {
		logger.Log.Error(fmt.Sprintf("Not supporting image extension %s", extension))
		return "", customError.NewServiceError(http.StatusBadRequest, fmt.Sprintf("Extension of image %s not support, pls provide correct one", extension), nil)
	}

	filename, err := fs.SafeImageToStorage(extension, &req.Image)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("Image saving error: %s", err.Error()))
		return "", customError.NewServiceError(http.StatusBadRequest, "Something went wrong while saving image. Check provided data or try later...", nil)
	}

	query, values := c.constructQueryForImageUpdate(req.Id, filename)
	ok, err = c.CatalogRepository.Update(cxt, query, values)

	if err != nil {
		logger.Log.Error("CatalogService -> AddImage -> err -> " + err.Error())
		return "", customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	if !ok {
		logger.Log.Error(fmt.Sprintf("CatalogService -> AddImage -> err -> "+"Cant update image of catalog item id: %d", req.Id))
		return "", customError.NewServiceError(http.StatusBadRequest, customError.ErrorNotFoundById, nil)
	}

	return filename, nil
}

func (c *CatalogService) GetAll(cxt context.Context, imagePrefix string) ([]*dto.CatalogResponse, error) {

	catalogItems, err := c.CatalogRepository.FindAll(cxt)

	if err != nil {
		logger.Log.Error("CatalogService -> GetAll -> err -> " + err.Error())
		return nil, customError.NewServiceError(http.StatusInternalServerError, customError.Error500, nil)
	}

	var catalogResp []*dto.CatalogResponse

	for _, item := range catalogItems {
		catalogResp = append(catalogResp, item.ToResponse(imagePrefix))
	}

	return catalogResp, nil
}

func (c *CatalogService) constructQueryForImageUpdate(userId uint, imageUrl string) (string, []any) {
	values := make([]any, 0)
	values = append(values, userId)

	query := fmt.Sprintf("image_url = $%d", 2)
	values = append(values, imageUrl)

	return query, values
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

func (c *CatalogService) getCatalogUniqFieldError(err error, catalog *model.Catalog) error {
	if strings.Contains(err.Error(), "sku") {
		return customError.NewServiceError(http.StatusConflict, fmt.Sprintf("Provided sku: %s already exist", catalog.Sku), err)
	}

	return customError.NewServiceError(http.StatusConflict, "Please check provided data", err)
}
