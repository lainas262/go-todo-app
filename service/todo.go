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
			for _, val := range users {
				var collaborator api.ResponseCollaborator
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

func (s *TodoService) DeleteTodo(ctx context.Context, userId int64, todoId int64) error {
	log.Printf("TodoService.DeleteTodo - Starting todo update for user:")

	err := s.repository.DeleteTodo(ctx, userId, todoId)

	if err != nil {
		log.Printf("TodoService.DeleteTodo - Database error: %v", err)
		return fmt.Errorf("failed to delete todo")
	}

	log.Printf("TodoService.DeleteTodo - Todo deleted successfully in database: %v", todoId)

	return nil
}

func (s *TodoService) UpdateCollaborators(ctx context.Context, authorId int64, todoId int64, req api.RequestAddCollaborators) error {

	log.Printf("TodoService.UpdateCollaborators - Starting todo update collaborators:")

	err := s.repository.UpdateCollaborators(ctx, authorId, todoId, req.Collaborators)

	if err != nil {
		log.Printf("TodoService.UpdateCollaborators - Database error: %v", err)
		return fmt.Errorf("failed to update collaborators")
	}

	log.Printf("TodoService.UpdateCollaborators - Todo collaborators updated successfully")

	return nil
}

func (s *TodoService) DeleteCollaborator(ctx context.Context, collaboratorEmail string, authorId int64, todoId int64) error {

	log.Printf("TodoService.DeleteCollaborator - Starting todo update collaborators:")

	err := s.repository.DeleteCollaborator(ctx, collaboratorEmail, authorId, todoId)

	if err != nil {
		log.Printf("TodoService.DeleteCollaborator - Database error: %v", err)
		return fmt.Errorf("failed to delete collaborators")
	}

	log.Printf("TodoService.DeleteCollaborator - Todo collaborator deleted successfully")

	return nil
}

func (s *TodoService) GetUserCollaborationTodos(ctx context.Context, userId int64) (*api.ResponseCollaborations, error) {
	log.Printf("TodoService.GetUserCollaborationTodos - Starting collaboration todos fetching for user:")

	todos, err := s.repository.GetUserCollaborationTodos(ctx, userId)

	if err != nil {
		log.Printf("TodoService.GetUserCollaborationTodos - Database error collaboration todo fetching: %v", err)
		return nil, fmt.Errorf("failed to fetch collaboration todos")
	}

	log.Printf("TodoService.GetUserCollaborationTodos - Collaboration todos fetched successfully")

	var todoIds []int64
	for _, val := range *todos {
		todoIds = append(todoIds, val.Id)
	}

	log.Printf("TodoService.GetUserCollaborationTodos - Starting collaborators fetching for collaboration todos:")
	collaborators, err := s.repository.GetTodoCollaboratosWithAuthor(ctx, todoIds)

	if err != nil {
		log.Printf("TodoService.GetUserCollaborationTodos - Database error collaborators fetching: %v", err)
		return nil, fmt.Errorf("failed to fetch todos")
	}

	log.Printf("TodoService.GetUserCollaborationTodos - Collaborators fetched successfully")

	var data []api.ResponseCollaborationTodo
	for _, todoModel := range *todos {
		var todoRes api.ResponseCollaborationTodo

		todoRes.Id = todoModel.Id
		todoRes.Title = todoModel.Title
		todoRes.Body = todoModel.Body
		todoRes.UpdatedAt = todoModel.UpdatedAt
		todoRes.CreatedAt = todoModel.CreatedAt

		users, ok := collaborators[todoModel.Id]
		if ok {
			for _, userModel := range users {
				var collaborator api.ResponseCollaborator
				collaborator.Email = userModel.Email
				collaborator.FirstName = userModel.FirstName
				collaborator.LastName = userModel.LastName
				collaborator.UserName = userModel.UserName

				if userModel.Id == userId {
					continue
				}

				if userModel.Id == todoModel.Author {
					todoRes.Author = collaborator
				} else {
					todoRes.Collaborators = append(todoRes.Collaborators, collaborator)
				}

			}
		}

		data = append(data, todoRes)
	}

	return &api.ResponseCollaborations{Data: data}, nil
}

func (s *TodoService) UpdateCollaborationTodo(ctx context.Context, userId int64, todoId int64, req api.RequestCreateTodo) (*api.ResponseCreateTodo, error) {
	log.Printf("TodoService.UpdateCollaborationTodo - Starting collaboration todo update for user:")

	if len(req.Title) > 255 {
		log.Printf("TodoService.UpdateCollaborationTodo - Validation failed: title length")
		return nil, fmt.Errorf("title must be up to 255 characters")
	}

	t := &model.Todo{
		Id:        todoId,
		Title:     req.Title,
		Body:      req.Body,
		Author:    userId,
		UpdatedAt: time.Now(),
	}

	todo, err := s.repository.UpdateCollaborationTodo(ctx, userId, t)

	if err != nil {
		log.Printf("TodoService.UpdateCollaborationTodo - Database error: %v", err)
		return nil, fmt.Errorf("failed to update todo")
	}

	log.Printf("TodoService.UpdateCollaborationTodo - Collaboration todo updated successfully in database: %v", todo.Id)

	return &api.ResponseCreateTodo{
		Id:        todo.Id,
		Title:     todo.Title,
		Body:      todo.Body,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}, nil
}
