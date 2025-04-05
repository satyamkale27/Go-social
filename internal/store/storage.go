package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage struct {
	// note: this is blueprint of NewStorage

	Posts interface {
		GetById(context.Context, int64) (*Post, error)
		Create(context.Context, *Post) error
	}
	users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db},
		users: &UsersStore{db},
	}
}

/*
personal note imp:-

The PostsStore struct has a method Create defined on it.
When you initialize the Posts field with &PostsStore{db},
you are creating a pointer to a PostsStore instance.
This pointer has access to all the methods defined on PostsStore,
including the Create method.

(PostsStore struct has a method Create defined on it,
code :  func (s *PostsStore) Create)
*/
