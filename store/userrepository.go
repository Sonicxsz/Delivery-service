package store

import (
	"arabic/internal/model"
	"context"
)

type UserRepository struct {
	store *Store
}

var (
	insertUser        = "INSERT INTO test.users (email, username, password) VALUES ($1, $2, $3) RETURNING id"
	searchUserByEmail = "SELECT id, username, password, email FROM test.users WHERE email = $1"
)

func (ur *UserRepository) Create(cxt context.Context, u *model.User) (*model.User, error) {
	err := ur.store.db.QueryRow(cxt, insertUser, u.Email, u.Username, u.Password).Scan(&u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) FindByEmail(cxt context.Context, email string) (*model.User, bool, error) {
	user := model.User{}
	err := ur.store.db.QueryRow(cxt, searchUserByEmail, email).Scan(&user.Id, &user.Username, &user.Password, &user.Email)

	if err != nil {
		return nil, false, err
	}

	return &user, true, nil
}
