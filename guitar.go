package main

import (
	"time"

	"github.com/google/uuid"
)

type Guitar struct {
	ID        uuid.UUID `json:"id"`
	Brand     string    `json:"brand"`
	Model     string    `json:"model"`
	Strings   int       `json:"strings"`
	CreatedAt time.Time `json:"created_at"`
}
