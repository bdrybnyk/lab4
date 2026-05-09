package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type mockGuitarRepo struct {
	guitars map[uuid.UUID]*Guitar
}

func newMockGuitarRepo() *mockGuitarRepo {
	return &mockGuitarRepo{guitars: make(map[uuid.UUID]*Guitar)}
}

func (m *mockGuitarRepo) GetByID(ctx context.Context, id uuid.UUID) (*Guitar, error) {
	g, ok := m.guitars[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return g, nil
}

func (m *mockGuitarRepo) Create(ctx context.Context, g *Guitar) error {
	m.guitars[g.ID] = g
	return nil
}

func (m *mockGuitarRepo) List(ctx context.Context, limit int, offset int) ([]Guitar, error) {
	var list []Guitar
	for _, g := range m.guitars {
		list = append(list, *g)
	}
	return list, nil
}

func (m *mockGuitarRepo) UpdateFull(ctx context.Context, g *Guitar) error {
	m.guitars[g.ID] = g
	return nil
}

func (m *mockGuitarRepo) UpdatePartial(ctx context.Context, id uuid.UUID, brand string) error {
	g, ok := m.guitars[id]
	if ok {
		g.Brand = brand
	}
	return nil
}

func (m *mockGuitarRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.guitars, id)
	return nil
}

func TestCreateHandler_InvalidJSON(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)

	req := httptest.NewRequest(http.MethodPost, "/guitars", strings.NewReader("{invalid}"))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetByIdHandler_NotFound(t *testing.T) {
	repo := newMockGuitarRepo()
	handler := NewGuitarHandler(repo)

	randomID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/guitars/"+randomID, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", randomID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler.GetById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
