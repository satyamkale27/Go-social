package store

import (
	"context"
	"database/sql"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
 INSERT INTO users (user_id, username, password, email) values ($1,$2,$3) RETURNING id, created_at
`
	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		return err

	}
	return nil
}
