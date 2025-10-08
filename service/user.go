package service

import (
	"context"
	"example/todolist/api"
	"example/todolist/auth"
	"example/todolist/model"
	"example/todolist/repository"
	"fmt"
	"log"
	"strings"
)

type UserService struct {
	repository *repository.UserRepository
}

func CreateUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(ctx context.Context, req api.RequestCreateUser) (*api.ResponseLoginUser, error) {

	log.Printf("UserService.CreateUser - Starting user creation for: %s", req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Printf("UserService.CreateUser - Validation failed: missing required fields")
		return nil, fmt.Errorf("username, email, and password are required")
	}

	if len(req.Password) < 8 {
		log.Printf("UserService.CreateUser - Validation failed: password length")
		return nil, fmt.Errorf("password must be at least 8 characters")
	}

	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		log.Printf("UserService.CreateUser - Error hashing password: %v", err)
		return nil, fmt.Errorf("error processing password")
	}

	u := &model.User{
		Email:        req.Email,
		UserName:     req.Username,
		PasswordHash: &hashedPassword,
	}

	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}

	if req.LastName != "" {
		u.LastName = req.LastName
	}

	user, err := s.repository.CreateUser(ctx, u)

	if err != nil {
		log.Printf("UserService.CreateUser - Database error: %v", err)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("username or email already exists")
		}
		return nil, fmt.Errorf("failed to create user")
	}

	log.Printf("UserService.CreateUser - User created successfully in database: %v", user.Id)

	return &api.ResponseLoginUser{
		AccessToken: "Lakis Very Secret Token",
		Username:    user.UserName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}, nil
}
