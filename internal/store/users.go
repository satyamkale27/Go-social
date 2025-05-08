package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail    = errors.New("a user with this email already exists")
	ErrDuplicateUsername = errors.New("a user with this username already exists")
)

type User struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
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

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
 INSERT INTO users ( username, password, email) values ($1,$2,$3) RETURNING id, created_at
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Password.hash, user.Email).Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violet unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violet unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err

		}

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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		if err := s.createUserninvitation(ctx, tx, token, invitationExp, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserninvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	// 1 find the user that this token belongs to
	// 2 update the user state(Active state)
	// 3 clean the invitations
	// so we need to run all the methods database operations, every operations has to be success so we use transactions
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		if err := s.deleteUserInvitations(ctx, tx, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {

	// update the user state(Active state)

	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active FROM users u JOIN user_invitations ui ON u.id = ui.user_id where ui.token = $1 AND ui.expiry > $2`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	var user User
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(

		&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err

		}
	}
	return &user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {

	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.Id)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userId int64) error {
	// clean the invitations
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}
	return nil
}
