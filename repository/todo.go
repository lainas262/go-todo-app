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
	SELECT tc.todo_id, u.email, u.first_name, u.last_name, u.username, u.id
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

		if err := rows.Scan(&todoId, &user.Email, &user.FirstName, &user.LastName, &user.UserName, &user.Id); err != nil {
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

func (r *TodoRepository) DeleteTodo(ctx context.Context, userId int64, todoId int64) error {
	query := `
	DELETE FROM todos
	WHERE author_id=$1 AND id=$2
	`
	cmdTag, err := r.db.Exec(ctx, query, userId, todoId)

	if err != nil {
		return fmt.Errorf("query delete todo error: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("query delete todo not found")
	}

	if cmdTag.RowsAffected() > 1 {
		return fmt.Errorf("warning: expected 1 row affected, got %d", cmdTag.RowsAffected())
	}
	return nil
}

func (r *TodoRepository) UpdateCollaborators(ctx context.Context, authorId int64, todoId int64, collaborators []string) error {
	query := `
	INSERT INTO todo_collaborators (user_id, todo_id)
	SELECT users.id, $2
	FROM users
	WHERE EXISTS (
		SELECT 1 FROM todos
		WHERE author_id=$1 AND id=$2
	)
	AND users.email = ANY($3)
	ON CONFLICT DO NOTHING
	`
	_, err := r.db.Query(ctx, query, authorId, todoId, collaborators)

	if err != nil {
		return fmt.Errorf("query update collaborators error: %w", err)
	}

	return nil
}

func (r *TodoRepository) DeleteCollaborator(ctx context.Context, collaboratorEmail string, authorId int64, todoId int64) error {
	query := `
	DELETE FROM todo_collaborators
	USING users, todos
	WHERE todos.id = todo_collaborators.todo_id
		AND users.id = todo_collaborators.user_id
		AND todos.id=$3
		AND todos.author_id=$2
		AND users.email=$1
	`

	_, err := r.db.Query(ctx, query, collaboratorEmail, authorId, todoId)

	if err != nil {
		return fmt.Errorf("query delete collaborator error: %w", err)
	}

	return nil
}

func (r *TodoRepository) GetUserCollaborationTodos(ctx context.Context, userId int64) (*[]model.Todo, error) {
	query := `
	SELECT t.id, t.title, t.body, t.created_at, t.updated_at, t.author_id 
	FROM todo_collaborators tc
	INNER JOIN todos t
	ON t.id = tc.todo_id
	WHERE tc.user_id=$1
	`
	rows, err := r.db.Query(ctx, query, userId)

	if err != nil {
		return nil, fmt.Errorf("query collaboration todos by user id: %w", err)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])

	if err != nil {
		return nil, fmt.Errorf("error collecting collaboration todo rows: %w", err)
	}

	return &todos, nil
}

func (r *TodoRepository) GetTodoCollaboratosWithAuthor(ctx context.Context, todoIds []int64) (map[int64][]model.User, error) {
	query := `
	SELECT t.id AS todo_id, u.email, u.first_name, u.last_name, u.username, u.id
	FROM users u
	LEFT JOIN todos t
	ON t.author_id=u.id
	WHERE t.id = ANY($1)
	UNION ALL
	SELECT tc.todo_id AS todo_id, u.email, u.first_name, u.last_name, u.username, u.id
	FROM users u
	LEFT JOIN todo_collaborators tc
	ON tc.user_id=u.id
	WHERE tc.todo_id = ANY($1)
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

		if err := rows.Scan(&todoId, &user.Email, &user.FirstName, &user.LastName, &user.UserName, &user.Id); err != nil {
			return nil, fmt.Errorf("error scanning todo collaborators join: %w", err)
		}
		collaboratorsMap[todoId] = append(collaboratorsMap[todoId], user)
	}

	return collaboratorsMap, nil

}

func (r *TodoRepository) UpdateCollaborationTodo(ctx context.Context, collaboratorId int64, todo *model.Todo) (*model.Todo, error) {
	query := `
	UPDATE todos
	SET title=$1, body=$2, updated_at=$3
	WHERE EXISTS (
		SELECT 1 FROM todo_collaborators
		WHERE user_id=$5 AND todo_id=$4
	)
	AND id=$4
	RETURNING created_at
	`
	err := r.db.QueryRow(ctx, query, todo.Title, todo.Body, todo.UpdatedAt, todo.Id, collaboratorId).Scan(&todo.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("query update todo error: %w", err)
	}
	return todo, nil
}
