package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
