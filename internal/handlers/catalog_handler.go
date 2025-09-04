package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"arabic/pkg/errors"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	id, ok := vars["id"]

	if id == "" || !ok {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, "Item Id not provided", nil), "Catalog: Delete item")
		return
	}

	parsedId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, "Cannot parse provided id", nil), "Catalog: Delete item")
		return
	}

	err = c.service.Delete(r.Context(), parsedId)

	if err != nil {
		handleServiceError(w, err, "Catalog: Delete item")
		return
	}

	respondSuccess(w, http.StatusOK, nil)

}

func (c *CatalogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.GetAll(r.Context())

	if err != nil {
		handleServiceError(w, err, "CategoryHandle GetManyById")
		return
	}
	respondSuccess(w, http.StatusOK, items)
}

func (c *CatalogHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemId, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, errors.ErrorGetQueryParam, nil), "CategoryHandle GetById")
		return
	}
	item, err := c.service.GetById(r.Context(), itemId)

	if err != nil {
		handleServiceError(w, err, "CategoryHandle GetManyById")
		return
	}

	respondSuccess(w, http.StatusOK, item)
	return
}

func (c *CatalogHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := dto.CatalogUpdateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, errors.ErrorParse, nil), "CategoryHandle Update")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, errors.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "CategoryHandle GetAll")
		return
	}

	err = c.service.Update(r.Context(), &req)

	if err != nil {
		handleServiceError(w, err, "CategoryHandle Update")
		return
	}

	respondSuccess(w, http.StatusOK, nil)
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

	id, err := c.service.Create(r.Context(), &req)

	if err != nil {
		err = errors.NewServiceError(http.StatusBadRequest, err.Error(), err)
		handleServiceError(w, err, "CategoryHandler Create")
		return
	}

	respondSuccess(w, http.StatusCreated, fmt.Sprintf("Id: %d", id))
}
