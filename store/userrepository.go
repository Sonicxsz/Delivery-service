package store

import (
	"arabic/internal/model"
	"context"
	"fmt"
)

type UserRepository struct {
	store *Store
}

var (
	table = "users"
)

func (ur *UserRepository) Create(u *model.User) (*model.User, error) {
	query := fmt.Sprintf("INSERT INTO %s (email, name, password) VALUES ($1, $2, $3) RETURNING id", table)
	err := ur.store.db.QueryRow(context.Background(), query, u.Email, u.Name, u.Password).Scan(&u.Id)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FindByEmail() {

}
