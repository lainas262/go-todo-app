package repository

import (
	"context"
	"example/todolist/model"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	db *pgxpool.Pool
}

func CreateTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetUserTodos(ctx context.Context, userId int64) (*[]model.Todo, error) {

	query := `
	SELECT id, title, body, created_at, updated_at, author_id
	FROM todos
	WHERE author_id=$1
	`

	rows, err := r.db.Query(ctx, query, userId)

	if err != nil {
		return nil, fmt.Errorf("query todos by author id: %w", err)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])

	if err != nil {
		return nil, fmt.Errorf("error collecting todo rows: %w", err)
	}

	return &todos, nil
}

func (r *TodoRepository) GetTodoCollaboratos(ctx context.Context, todoIds []int64) (map[int64][]model.User, error) {
	query := `
	SELECT tc.todo_id, u.email, u.first_name, u.last_name, u.username
	FROM todo_collaborators tc
	INNER JOIN users u
	ON tc.user_id=u.id
	WHERE todo_id = ANY($1)
	`

	rows, err := r.db.Query(ctx, query, todoIds)

	if err != nil {
		return nil, fmt.Errorf("query todo collaborators join: %w", err)
	}

	defer rows.Close()

	collaboratorsMap := make(map[int64][]model.User)
	for rows.Next() {
		var todoId int64
		var user model.User

		if err := rows.Scan(&todoId, &user.Email, &user.FirstName, &user.LastName, &user.UserName); err != nil {
			return nil, fmt.Errorf("error scanning todo collaborators join: %w", err)
		}
		collaboratorsMap[todoId] = append(collaboratorsMap[todoId], user)
	}

	return collaboratorsMap, nil

}

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	query := `
	INSERT INTO todos (author_id, title, body)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, todo.Author, todo.Title, todo.Body).Scan(&todo.Id, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("query create todo error: %w", err)
	}

	return todo, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	query := `
	UPDATE todos
	SET title=$1, body=$2, updated_at=$3
	WHERE id=$4 AND author_id=$5
	RETURNING created_at
	`
	err := r.db.QueryRow(ctx, query, todo.Title, todo.Body, todo.UpdatedAt, todo.Id, todo.Author).Scan(&todo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("query update todo error: %w", err)
	}
	return todo, nil
}
