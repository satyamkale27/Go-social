package store

import (
	"context"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
 INSERT INTO users ( username, password, email) values ($1,$2,$3) RETURNING id, created_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		return err

	}
	return nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {

	query := `SELECT id, email, username, password,  created_at
             FROM users
             WHERE id = $1
             `
	var user User

	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id, &user.Email, &user.Username, &user.Password, &user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string) error {

}
