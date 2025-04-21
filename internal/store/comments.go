package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"` // User is a struct

}

type comentStore struct {
	db *sql.DB
}

func (s *comentStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
               SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id  FROM comments c
              JOIN users ON users.id = c.user_id
              WHERE c.post_id= $1
              ORDER BY c.created_at DESC;
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		/*
			 c.User = User{}
			This explicitly initializes the User field of the Comment struct to an empty User struct.
		*/

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, c.User.Username, c.User.Id)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
		/*
			The append function in Go is used to add elements to a slice.
			It takes a slice and one or more elements as arguments and returns a new slice
			with the elements added.
		*/
	}
	return comments, nil
}

/*
 []Comment{} is a slice of Comment structs.
It can hold multiple Comment structs,
allowing you to store and manage a collection of comments in a single variable.
Each element in the slice is of type Comment.
*/

func (s *comentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id,created_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
