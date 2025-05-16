package store

import (
	"context"
	"database/sql"
)

type Role struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       int64  `json:"level"`
}

type RoleStore struct {
	db *sql.DB
}

func (r *RoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`
	role := &Role{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(&role.Id, &role.Name, &role.Description, &role.Level)

	if err != nil {

		return nil, err
	}
	return role, nil
}
