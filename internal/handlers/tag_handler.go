package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TagHandler struct {
	service service.ITagService
}

func NewTagHandler(service service.ITagService) *TagHandler {
	return &TagHandler{service: service}
}

func (t *TagHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := t.service.GetAll(r.Context())
		if err != nil {
			handleServiceError(w, err, "TagHandle GetAll")
			return
		}
		respondSuccess(w, http.StatusOK, all)
	}
}

func (t *TagHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.TagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		tag, err := t.service.Create(r.Context(), &req)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Something went wrong. pls try later")
			return
		}
		respondSuccess(w, http.StatusOK, tag)
	}
}

func (t *TagHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		tagId, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid tag ID")
			return
		}

		err = t.service.Delete(r.Context(), tagId)
		if err != nil {
			handleServiceError(w, err, "TagHandle Delete")
			return
		}

		respondSuccess(w, http.StatusOK, nil)
	}
}
