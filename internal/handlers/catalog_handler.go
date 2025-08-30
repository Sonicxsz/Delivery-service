package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"arabic/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type CatalogHandler struct {
	service service.ICatalogService
}

func NewCatalogHandler(service service.ICatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

func (c *CatalogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, "Item Id not provided", nil), "Catalog: Delete item")
		return
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, "Cannot parse provided id", nil), "Catalog: Delete item")
		return
	}

	err = c.service.Delete(context.Background(), parsedId)

	if err != nil {
		handleServiceError(w, err, "Catalog: Delete item")
		return
	}

	respondSuccess(w, http.StatusOK, nil)

}

func (c *CatalogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.GetAll(context.Background())

	if err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusInternalServerError, errors.Error500, nil), "Catalog: Get All Catalog Items")
	}
	respondSuccess(w, http.StatusOK, items)
}

func (c *CatalogHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := dto.CatalogCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, errors.ErrorParse, nil), "Catalog: Create Decode")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "CategoryHandle GetAll")
		return
	}

	id, err := c.service.Create(context.Background(), &req)

	if err != nil {
		err = errors.NewServiceError(http.StatusBadRequest, err.Error(), err)
		handleServiceError(w, err, "CategoryHandler Create")
		return
	}

	respondSuccess(w, http.StatusCreated, fmt.Sprintf("Id: %d", id))
}
