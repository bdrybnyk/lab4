package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GuitarRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Guitar, error)
	Create(ctx context.Context, g *Guitar) error
	List(ctx context.Context, limit int, offset int) ([]Guitar, error)
	UpdateFull(ctx context.Context, g *Guitar) error
	UpdatePartial(ctx context.Context, id uuid.UUID, brand string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type guitarRepo struct {
	db *pgxpool.Pool
}

func NewGuitarRepo(db *pgxpool.Pool) GuitarRepo {
	return &guitarRepo{db: db}
}

func (r *guitarRepo) GetByID(ctx context.Context, id uuid.UUID) (*Guitar, error) {
	query := `SELECT id, brand, model, strings, created_at FROM guitar WHERE id = $1`
	g := &Guitar{}
	err := r.db.QueryRow(ctx, query, id).Scan(&g.ID, &g.Brand, &g.Model, &g.Strings, &g.CreatedAt)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (r *guitarRepo) Create(ctx context.Context, g *Guitar) error {
	query := `INSERT INTO guitar (id, brand, model, strings) VALUES ($1, $2, $3, $4) RETURNING created_at`
	err := r.db.QueryRow(ctx, query, g.ID, g.Brand, g.Model, g.Strings).Scan(&g.CreatedAt)
	return err
}

func (r *guitarRepo) List(ctx context.Context, limit int, offset int) ([]Guitar, error) {
	query := `SELECT id, brand, model, strings, created_at FROM guitar LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guitars []Guitar
	for rows.Next() {
		var g Guitar
		err := rows.Scan(&g.ID, &g.Brand, &g.Model, &g.Strings, &g.CreatedAt)
		if err != nil {
			return nil, err
		}
		guitars = append(guitars, g)
	}
	return guitars, nil
}

func (r *guitarRepo) UpdateFull(ctx context.Context, g *Guitar) error {
	query := `UPDATE guitar SET brand = $1, model = $2, strings = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, g.Brand, g.Model, g.Strings, g.ID)
	return err
}

func (r *guitarRepo) UpdatePartial(ctx context.Context, id uuid.UUID, brand string) error {
	query := `UPDATE guitar SET brand = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, brand, id)
	return err
}

func (r *guitarRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM guitar WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
