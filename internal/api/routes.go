package api

import (
	"log"
	"net/http"

	"github.com/shenikar/Name-analyzer/internal/db"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/shenikar/Name-analyzer/docs"
)

func RegisterRoutes(mux *http.ServeMux, database *db.DB, logger *log.Logger) {
	h := &Handler{
		DB:     database,
		Logger: logger,
	}
	mux.HandleFunc("POST /api/v1/persons", h.CreatePerson)
	mux.HandleFunc("GET /api/v1/persons", h.ListPersons)
	mux.HandleFunc("GET /api/v1/persons/{id}", h.GetPerson)
	mux.HandleFunc("PUT /api/v1/persons/{id}", h.UpdatePerson)
	mux.HandleFunc("DELETE /api/v1/persons/{id}", h.DeletePerson)

	// Swagger UI
    mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
}
