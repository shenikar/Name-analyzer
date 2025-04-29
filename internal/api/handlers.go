package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/shenikar/Name-analyzer/internal/db"
	"github.com/shenikar/Name-analyzer/internal/enrich"
	"github.com/shenikar/Name-analyzer/internal/model"
)

type Handler struct {
	DB     *db.DB
	Logger *log.Logger
}

var req struct {
	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Age         *int    `json:"age"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`
}

// CreatePerson godoc
// @Summary Создать новую запись о человеке
// @Description Создает новую запись и обогащает её данными о возрасте, поле и национальности
// @Tags persons
// @Accept json
// @Produce json
// @Param request body model.PersonRequest true "Данные о человеке"
// @Success 201 {object} model.Person
// @Failure 400 {object} model.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [post]
func (h *Handler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Surname = strings.TrimSpace(req.Surname)
	if req.Name == "" || req.Surname == "" {
		http.Error(w, "name and surname are required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	data, _ := enrich.EnrichPerson(ctx, req.Name)

	person := &model.Person{
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		Age:         data.Age,
		Gender:      data.Gender,
		Nationality: data.Nationality,
	}
	if err := h.DB.CreatePerson(ctx, person); err != nil {
		h.Logger.Printf("failed to create person: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// GetPerson godoc
// @Summary Получить информацию о человеке по ID
// @Description Возвращает детальную информацию о человеке
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека" format(uuid)
// @Success 200 {object} model.Person
// @Failure 400 {object} model.ErrorResponse "Некорректный ID"
// @Failure 404 {object} model.ErrorResponse "Человек не найден"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [get]
func (h *Handler) GetPerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/persons/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	person, err := h.DB.GetPerson(r.Context(), id)
	if err == db.ErrNotFound {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Logger.Printf("failed to get person: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// ListPersons godoc
// @Summary Получить список людей
// @Description Возвращает список людей с возможностью фильтрации
// @Tags persons
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param gender query string false "Фильтр по полу"
// @Param nationality query string false "Фильтр по национальности"
// @Param age_min query integer false "Минимальный возраст"
// @Param age_max query integer false "Максимальный возраст"
// @Param limit query integer false "Количество записей на странице" default(10)
// @Param offset query integer false "Смещение" default(0)
// @Success 200 {array} model.Person
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons [get]
func (h *Handler) ListPersons(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := map[string]interface{}{}
	if v := q.Get("name"); v != "" {
		filter["name"] = v
	}
	if v := q.Get("surname"); v != "" {
		filter["surnamename"] = v
	}
	if v := q.Get("gender"); v != "" {
		filter["gender"] = v
	}
	if v := q.Get("nationality"); v != "" {
		filter["nationality"] = v
	}
	if v := q.Get("age_min"); v != "" {
		filter["age_min"] = v
	}
	if v := q.Get("age_max"); v != "" {
		filter["age_max"] = v
	}
	limit := 10
	offset := 0
	if v := q.Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &limit)
	}
	if v := q.Get("offset"); v != "" {
		fmt.Sscanf(v, "%d", &offset)
	}
	persons, err := h.DB.ListPersons(r.Context(), filter, limit, offset)
	if err != nil {
		h.Logger.Printf("error listing persons: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// UpdatePerson godoc
// @Summary Обновить информацию о человеке
// @Description Обновляет существующую запись о человеке
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека" format(uuid)
// @Param request body model.PersonRequest true "Обновленные данные"
// @Success 200 {object} model.Person
// @Failure 400 {object} model.ErrorResponse "Некорректный запрос"
// @Failure 404 {object} model.ErrorResponse "Человек не найден"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [put]
func (h *Handler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/persons/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	person, err := h.DB.GetPerson(r.Context(), id)
	if err == db.ErrNotFound {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Logger.Printf("failed to update person: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	if req.Name != "" {
		person.Name = req.Name
	}
	if req.Surname != "" {
		person.Surname = req.Surname
	}
	if req.Patronymic != nil {
		person.Patronymic = req.Patronymic
	}
	if req.Age != nil {
		person.Age = req.Age
	}
	if req.Gender != nil {
		person.Gender = req.Gender
	}
	if req.Nationality != nil {
		person.Nationality = req.Nationality
	}
	if err := h.DB.UpdatePerson(r.Context(), person); err != nil {
		h.Logger.Printf("failed to update person: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// DeletePerson godoc
// @Summary Удалить запись о человеке
// @Description Удаляет запись о человеке по ID
// @Tags persons
// @Accept json
// @Produce json
// @Param id path string true "ID человека" format(uuid)
// @Success 204 "Запись успешно удалена"
// @Failure 400 {object} model.ErrorResponse "Некорректный ID"
// @Failure 404 {object} model.ErrorResponse "Человек не найден"
// @Failure 500 {object} model.ErrorResponse "Внутренняя ошибка сервера"
// @Router /persons/{id} [delete]
func (h *Handler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "api/v1/persons/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = h.DB.DeletePerson(r.Context(), id)
	if err == db.ErrNotFound {
		http.Error(w, "person not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.Logger.Printf("failed to delete person: %v", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
