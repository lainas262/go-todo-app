package repository

import (
	"context"
	"errors"
	"example/todolist/model"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func CreateUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserById(ctx context.Context, id int64) (*model.User, error) {

	query := `
	SELECT id, first_name, last_name, email, username, password_hash
	FROM users
	WHERE id=$1
	`

	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.UserName, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found") // User not found
		}
		return nil, fmt.Errorf("query user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	query := `
	INSERT INTO users (first_name, last_name, email, username, password_hash) 
	VALUES ($1, $2, $3, $4, $5) 
	RETURNING id
	`

	err := r.db.QueryRow(context.Background(), query, user.FirstName, user.LastName, user.Email, user.UserName, user.PasswordHash).Scan(&user.Id)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, errors.New("email already exists")
			}
		}
		return nil, fmt.Errorf("inder user: %w", err)
	}

	return user, nil
}
