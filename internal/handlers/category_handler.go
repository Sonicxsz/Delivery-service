package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	service service.ICategoryService
}

func NewCategoryHandler(service service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (t *CategoryHandler) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := t.service.FindAll(r.Context())
		if err != nil {
			handleServiceError(w, err, "FindAll")
			return
		}
		respondSuccess(w, http.StatusOK, all)
	}
}

func (t *CategoryHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		category, err := t.service.Create(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Something went wrong. pls try later")
			return
		}
		respondSuccess(w, http.StatusOK, category)
	}
}

func (t *CategoryHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		categoryId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid category ID")
			return
		}

		err = t.service.Delete(r.Context(), categoryId)
		if err != nil {
			handleServiceError(w, err, "Delete")
			return
		}

		respondSuccess(w, http.StatusOK, nil)
	}
}
