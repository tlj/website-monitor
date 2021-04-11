package users

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Find(id int64) (*User, error) {
	u := &User{}

	if err := r.db.Get(u, "SELECT * FROM users WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("error looking up user id '%d': %s", id, err)
	}

	return u, nil
}
