package store

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Follower struct {
	UserId     int64 `db:"user_id"`
	FollowerId int64 `db:"follower_id"`
	CreatedAt  int64 `db:"created_at"`
}
type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerId, userId int64) error {

	query := `INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followerId, userId int64) error {

	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		return err
	}
	return nil
}
