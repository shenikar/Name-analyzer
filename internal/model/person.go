package model

import (
	"time"

	"github.com/google/uuid"
)

// Person представляет информацию о человеке с обогащенными данными
type Person struct {
	ID          uuid.UUID `db:"id" json:"id" example:"39755c70-2ddb-4a62-90ea-1eeaf07a545a"`
	Name        string    `db:"name" json:"name" example:"Иван"`
	Surname     string    `db:"surname" json:"surname" example:"Иванов"`
	Patronymic  *string   `db:"patronymic" json:"patronymic,omitempty" example:"Иванович"`
	Age         *int      `db:"age" json:"age,omitempty" example:"30"`
	Gender      *string   `db:"gender" json:"gender,omitempty" example:"male"`
	Nationality *string   `db:"nationality" json:"nationality,omitempty" example:"RU"`
	CreatedAt   time.Time `db:"created_at" json:"created_at" example:"2024-03-20T15:04:05Z"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at" example:"2024-03-20T15:04:05Z"`
}

// PersonRequest представляет запрос на создание/обновление записи
type PersonRequest struct {
	Name        string  `json:"name" example:"Иван"`
	Surname     string  `json:"surname" example:"Иванов"`
	Patronymic  *string `json:"patronymic,omitempty" example:"Иванович"`
	Age         *int    `json:"age,omitempty" example:"30"`
	Gender      *string `json:"gender,omitempty" example:"male"`
	Nationality *string `json:"nationality,omitempty" example:"RU"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"некорректный запрос"`
}
