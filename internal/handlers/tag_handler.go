package handlers

import (
	"arabic/internal/dto"
	"arabic/internal/service"
	"encoding/json"
	"net/http"
)

func FindAllTags(s *service.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := s.FindAll(r.Context())
		if err != nil {

		}
		respondSuccess(w, http.StatusOK, all)
	}
}

func CreateTag(s *service.TagService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req *dto.TagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		tag, err := s.CreateTag(r.Context(), req)
		if err != nil {

		}

		respondSuccess(w, http.StatusOK, tag)
	}
}
