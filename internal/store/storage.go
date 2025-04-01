package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
	}
	users interface {
		Create(context.Context) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
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
