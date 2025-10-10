package service

import (
	"context"
	"example/todolist/api"
	"example/todolist/model"
	"example/todolist/repository"
	"fmt"
	"log"
	"time"
)

type TodoService struct {
	repository *repository.TodoRepository
}

func CreateTodoService(repository *repository.TodoRepository) *TodoService {
	return &TodoService{repository: repository}
}

func (s *TodoService) GetUserTodos(ctx context.Context, userId int64) (*api.ResponseGetTodos, error) {
	log.Printf("TodoService.GetUserTodos - Starting todos fetching for user:")

	todos, err := s.repository.GetUserTodos(ctx, userId)

	if err != nil {
		log.Printf("TodoService.GetUserTodos - Database error todo fetching: %v", err)
		return nil, fmt.Errorf("failed to fetch todos")
	}

	log.Printf("TodoService.GetUserTodos - Todos fetched successfully")

	var todoIds []int64
	for _, val := range *todos {
		todoIds = append(todoIds, val.Id)
	}

	log.Printf("TodoService.GetUserTodos - Starting collaborators fetching for todos:")
	collaborators, err := s.repository.GetTodoCollaboratos(ctx, todoIds)

	if err != nil {
		log.Printf("TodoService.GetUserTodos - Database error collaborators fetching: %v", err)
		return nil, fmt.Errorf("failed to fetch todos")
	}

	log.Printf("TodoService.GetUserTodos - Collaborators fetched successfully")

	var data []api.ResponseTodoWithCollaborators
	for _, val := range *todos {
		var todoRes api.ResponseTodoWithCollaborators

		todoRes.Id = val.Id
		todoRes.Title = val.Title
		todoRes.Body = val.Body
		todoRes.UpdatedAt = val.UpdatedAt
		todoRes.CreatedAt = val.CreatedAt

		users, ok := collaborators[val.Id]
		if ok {
			log.Println("FOUND USERS ARRAY")

			for _, val := range users {
				var collaborator api.Collaborator
				collaborator.Email = val.Email
				collaborator.FirstName = val.FirstName
				collaborator.LastName = val.LastName
				collaborator.UserName = val.UserName

				todoRes.Collaborators = append(todoRes.Collaborators, collaborator)
			}
		}

		data = append(data, todoRes)
	}

	return &api.ResponseGetTodos{Data: data}, nil
}

func (s *TodoService) CreateTodo(ctx context.Context, userId int64, req api.RequestCreateTodo) (*api.ResponseCreateTodo, error) {

	log.Printf("TodoService.CreateTodo - Starting todo creation for user:")

	if req.Title == "" || req.Body == "" {
		log.Printf("TodoService.CreateTodo - Validation failed: missing required fields")
		return nil, fmt.Errorf("title and body are required")
	}

	if len(req.Title) > 255 {
		log.Printf("TodoService.CreateTodo - Validation failed: title length")
		return nil, fmt.Errorf("title must be up to 255 characters")
	}

	t := &model.Todo{
		Title:  req.Title,
		Body:   req.Body,
		Author: userId,
	}

	todo, err := s.repository.CreateTodo(ctx, t)

	if err != nil {
		log.Printf("TodoService.CreateTodo - Database error: %v", err)
		return nil, fmt.Errorf("failed to create todo")
	}

	log.Printf("TodoService.CreateTodo - Todo created successfully in database: %v", todo.Id)

	return &api.ResponseCreateTodo{
		Id:        todo.Id,
		Title:     todo.Title,
		Body:      todo.Body,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, userId int64, todoId int64, req api.RequestCreateTodo) (*api.ResponseCreateTodo, error) {
	log.Printf("TodoService.UpdateTodo - Starting todo update for user:")

	if len(req.Title) > 255 {
		log.Printf("TodoService.UpdateTodo - Validation failed: title length")
		return nil, fmt.Errorf("title must be up to 255 characters")
	}

	t := &model.Todo{
		Id:        todoId,
		Title:     req.Title,
		Body:      req.Body,
		Author:    userId,
		UpdatedAt: time.Now(),
	}

	todo, err := s.repository.UpdateTodo(ctx, t)

	if err != nil {
		log.Printf("TodoService.UpdateTodo - Database error: %v", err)
		return nil, fmt.Errorf("failed to update todo")
	}

	log.Printf("TodoService.UpdateTodo - Todo updated successfully in database: %v", todo.Id)

	return &api.ResponseCreateTodo{
		Id:        todo.Id,
		Title:     todo.Title,
		Body:      todo.Body,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}, nil
}
