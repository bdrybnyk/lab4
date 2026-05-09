package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type GuitarHandler struct {
	repo GuitarRepo
}

func NewGuitarHandler(repo GuitarRepo) *GuitarHandler {
	return &GuitarHandler{repo: repo}
}

func (h *GuitarHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	g, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		log.Printf("GetByID error: %v", err)
		http.Error(w, "guitar not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GuitarHandler) Create(w http.ResponseWriter, r *http.Request) {
	var g Guitar
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(g.Brand) == "" {
		http.Error(w, "brand must not be empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(g.Model) == "" {
		http.Error(w, "model must not be empty", http.StatusBadRequest)
		return
	}

	if g.Strings <= 0 {
		http.Error(w, "strings must be greater than zero", http.StatusBadRequest)
		return
	}

	g.ID = uuid.New()

	err = h.repo.Create(r.Context(), &g)
	if err != nil {
		log.Printf("failed to create guitar: %v\n", err)
		http.Error(w, "failed to create guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
}

func (h *GuitarHandler) List(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	guitars, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		log.Printf("failed to list guitars: %v\n", err)
		http.Error(w, "failed to list guitars", http.StatusInternalServerError)
		return
	}

	if guitars == nil {
		guitars = []Guitar{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitars)
}

func (h *GuitarHandler) UpdateFull(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var g Guitar
	err = json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(g.Brand) == "" {
		http.Error(w, "brand must not be empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(g.Model) == "" {
		http.Error(w, "model must not be empty", http.StatusBadRequest)
		return
	}

	if g.Strings <= 0 {
		http.Error(w, "strings must be greater than zero", http.StatusBadRequest)
		return
	}

	g.ID = id

	err = h.repo.UpdateFull(r.Context(), &g)
	if err != nil {
		log.Printf("failed to update guitar: %v\n", err)
		http.Error(w, "failed to update guitar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func (h *GuitarHandler) UpdatePartial(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var payload struct {
		Brand string `json:"brand"`
	}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(payload.Brand) == "" {
		http.Error(w, "brand must not be empty", http.StatusBadRequest)
		return
	}

	err = h.repo.UpdatePartial(r.Context(), id, payload.Brand)
	if err != nil {
		log.Printf("failed to partial update guitar: %v\n", err)
		http.Error(w, "failed to update guitar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GuitarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		log.Printf("failed to delete guitar: %v\n", err)
		http.Error(w, "failed to delete guitar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
