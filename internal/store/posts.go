package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	Id        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comment   []Comment `json:"comment"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {

	query := `
 INSERT INTO posts (content, title, user_id, tags) 
 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
 `
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {

	query := `SELECT id, content, title, user_id, tags, created_at, updated_at,version
             FROM posts
             WHERE id = $1
             `
	var post Post

	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.Id, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {

	query := `DELETE FROM posts WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {

	query := `UPDATE posts SET title = $1, content = $2, version = version + 1  WHERE id = $3 AND version = $4 RETURNING version`

	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.Id, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err

		}
	}
	return nil
}
