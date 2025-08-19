package store

import (
	"arabic/internal/model"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	store *Store
}

var (
	insertUser        = "INSERT INTO test.users (email, username, password) VALUES ($1, $2, $3) RETURNING id"
	searchUserByEmail = "SELECT id, username, password, email FROM test.users WHERE email = $1"
)

func (ur *UserRepository) Create(u *model.User) (string, error) {
	hashedPassword, err2 := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err2 != nil {
		return "", fmt.Errorf("failed to hash password: %w", err2)
	}

	err := ur.store.db.
		QueryRow(context.Background(), insertUser, u.Email, u.Username, string(hashedPassword)).
		Scan(&u.Id)

	if err != nil {
		return "", err
	}

	return u.Id, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, bool, error) {

	user := model.User{}
	err := ur.store.db.QueryRow(context.Background(), searchUserByEmail, email).Scan(&user.Id, &user.Username, &user.Password, &user.Email)

	if err != nil {
		return nil, false, err
	}

	return &user, true, nil
}
