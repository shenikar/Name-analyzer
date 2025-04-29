package db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shenikar/Name-analyzer/internal/model"
)

var ErrNotFound = errors.New("person not found")

func (db *DB) CreatePerson(ctx context.Context, person *model.Person) error {
	query := `
	     INSERT INTO persons (id, name, surname, patronymic, age, gender, nationality, created_at, updated_at)
		 VALUES (:id, :name, :surname, :patronymic, :age, :gender, :nationality, NOW(), NOW())
		 RETURNING created_at, updated_at
	`
	person.ID = uuid.New()
	rows, err := db.Conn.NamedQueryContext(ctx, query, person)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&person.CreatedAt, &person.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) GetPerson(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	var person model.Person
	err := db.Conn.GetContext(ctx, &person, `SELECT * FROM persons WHERE id=$1`, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &person, nil
}

func (db *DB) UpdatePerson(ctx context.Context, person *model.Person) error {
	query := `
        UPDATE persons SET name=:name, surname=:surname, patronymic=:patronymic, age=:age, gender=:gender, nationality=:nationality, updated_at=NOW()
		WHERE id=:id
		RETURNING updated_at
	`
	rows, err := db.Conn.NamedQueryContext(ctx, query, person)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&person.UpdatedAt); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) DeletePerson(ctx context.Context, id uuid.UUID) error {
	res, err := db.Conn.ExecContext(ctx, `DELETE FROM persons WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *DB) ListPersons(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*model.Person, error) {
	query := `SELECT * FROM persons WHERE 1=1`
	args := map[string]interface{}{}

	if v, ok := filter["name"]; ok {
		query += " AND name ILIKE :name"
		args["name"] = "%" + v.(string) + "%"
	}
	if v, ok := filter["surname"]; ok {
		query += " AND surname ILIKE :surname"
		args["surname"] = "%" + v.(string) + "%"
	}
	if v, ok := filter["gender"]; ok {
		query += " AND gender = :gender"
		args["gender"] = v
	}
	if v, ok := filter["nationality"]; ok {
		query += " AND nationality = :nationality"
		args["nationality"] = v
	}
	if v, ok := filter["age_min"]; ok {
		query += " AND age >= :age_min"
		args["age_min"] = v
	}
	if v, ok := filter["age_max"]; ok {
		query += " AND age <= :age_max"
		args["age_max"] = v
	}
	query += " ORDER BY created_at DESC LIMIT :limit OFFSET :offset"
	args["limit"] = limit
	args["offset"] = offset

	rows, err := db.Conn.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.StructScan(&person); err != nil {
			return nil, err
		}
		result = append(result, &person)
	}
	return result, nil

}
