package model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Surname     string    `db:"surname" json:"surname"`
	Patronymic  *string   `db:"patronymic" json:"patronymic,omitempty"`
	Age         *int      `db:"age" json:"age,omitempty"`
	Gender      *string   `db:"gender" json:"gender,omitempty"`
	Nationality *string   `db:"nationality" json:"nationality,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
