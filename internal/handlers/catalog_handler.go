package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"arabic/pkg/customError"
	"arabic/pkg/fs"
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

	if !ok {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, "Item Id not provided", nil), "Catalog: Delete item")
		return
	}

	parsedId, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, "Cannot parse provided id", nil), "Catalog: Delete item")
		return
	}

	err = c.service.Delete(r.Context(), uint(parsedId))

	if err != nil {
		handleServiceError(w, err, "Catalog: Delete item")
		return
	}

	respondSuccess(w, http.StatusOK, nil)

}

func (c *CatalogHandler) GetAll(fs fs.IFileSystemImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imagePrefix := "/" + fs.GetPath()
		items, err := c.service.GetAll(r.Context(), imagePrefix)

		if err != nil {
			handleServiceError(w, err, "CategoryHandle GetManyById")
			return
		}
		respondSuccess(w, http.StatusOK, items)
	}
}

func (c *CatalogHandler) GetById(fs fs.IFileSystemImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemId, err := strconv.ParseUint(vars["id"], 10, 0)

		if err != nil {
			handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorGetQueryParam, nil), "CategoryHandle GetById")
			return
		}

		imagePrefix := "/" + fs.GetPath()
		item, err := c.service.GetById(r.Context(), uint(itemId), imagePrefix)

		if err != nil {
			handleServiceError(w, err, "CategoryHandle GetManyById")
			return
		}

		respondSuccess(w, http.StatusOK, item)
	}
}

func (c *CatalogHandler) Update(w http.ResponseWriter, r *http.Request) {
	req := dto.CatalogUpdateRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorParse, nil), "Catalog: Parse error")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "Catalog: validation error")
		return
	}

	err = c.service.Update(r.Context(), &req)

	if err != nil {
		handleServiceError(w, err, "Catalog: Service error")
		return
	}

	respondSuccess(w, http.StatusOK, nil)
}

func (c *CatalogHandler) Create(w http.ResponseWriter, r *http.Request) {
	req := dto.CatalogCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorParse, nil), "Catalog: Create Decode")
		return
	}

	if ok, errStrings := req.IsValid(); !ok {
		handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, strings.Join(errStrings, "; "), nil), "CategoryHandle GetAll")
		return
	}

	id, err := c.service.Create(r.Context(), &req)

	if err != nil {
		err = customError.NewServiceError(http.StatusBadRequest, err.Error(), err)
		handleServiceError(w, err, "CategoryHandler Create")
		return
	}

	respondSuccess(w, http.StatusCreated, fmt.Sprintf("Id: %d", id))
}

func (c *CatalogHandler) AddImage(fs fs.IFileSystemImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &dto.AddImageRequest{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			handleServiceError(w, customError.NewServiceError(http.StatusBadRequest, customError.ErrorParse, nil), "Catalog: AddImage Decode")
			return
		}

		filename, err := c.service.AddImage(r.Context(), req, fs)

		if err != nil {
			handleServiceError(w, err, "Catalog: AddImage Service")
			return
		}

		respondSuccess(w, 200, filename)
	}
}
