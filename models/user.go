package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" validate:"uuid"`
	Name      string    `json:"name" validate:"required,min=2,max=50"`
	Age       int       `json:"age" validate:"gte=0,lte=120"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}
