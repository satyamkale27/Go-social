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
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post               // it is post embedding
	CommentCount int64 `json:"comment_count"`
}
type AllUserPosts struct {
	Post
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userid int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.content,
    p.created_at,
    p.version,
    p.tags,
    u.username,
    COUNT(c.id) AS comments_count
FROM 
    posts p
LEFT JOIN 
    comments c ON c.post_id = p.id
LEFT JOIN 
    users u ON p.user_id = u.id
JOIN 
    followers f ON f.follower_id = p.user_id OR p.user_id = $1
WHERE 
    f.user_id = $1  AND
    (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%' ) AND 
    (p.tags @> $5 OR $5 = '{}'  )
GROUP BY 
    p.id, u.username
ORDER BY 
    p.created_at ` + fq.Sort + `
LIMIT $2 OFFSET $3;
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userid, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var feed []PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err := rows.Scan(&p.Id, &p.UserID, &p.Title, &p.Content, &p.CreatedAt, &p.Version, pq.Array(&p.Tags), &p.User.Username, &p.CommentCount)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}
	return feed, nil
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {

	query := `
 INSERT INTO posts (content, title, user_id, tags) 
 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
 `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.Id, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `SELECT p.id, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at, p.version
             FROM posts p
             WHERE p.id = $1`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

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

func (s *PostStore) GetAllUserPosts(ctx context.Context, userid int64) ([]AllUserPosts, error) {
	query := `SELECT p.id,p.title,p.content,p.created_at,p.version,p.tags FROM posts p JOIN users ON p.user_id = users.id WHERE user_id = $1 ORDER BY p.created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userPosts []AllUserPosts
	for rows.Next() {
		var up AllUserPosts
		err := rows.Scan(&up.Id, &up.Title, &up.Content, &up.CreatedAt, &up.Version, pq.Array(&up.Tags))
		if err != nil {
			return nil, err
		}
		userPosts = append(userPosts, up)
	}

	return userPosts, nil

}
