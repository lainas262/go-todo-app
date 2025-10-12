package api

import "time"

type RequestCreateTodo struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type ResponseCreateTodo struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RequestAddCollaborators struct {
	Collaborators []string `json:"collaborators"`
}

type RequestDeleteCollaborator struct {
	Collaborator string `json:"collaborator"`
}

type ResponseCollaborator struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type ResponseTodoWithCollaborators struct {
	Id            int64                  `json:"id"`
	Title         string                 `json:"title"`
	Body          string                 `json:"body"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CreatedAt     time.Time              `json:"created_at"`
	Collaborators []ResponseCollaborator `json:"collaborators"`
}
type ResponseGetTodos struct {
	Data []ResponseTodoWithCollaborators `json:"data"`
}

// essentially enhanced ResponseTodoWithCollaborators which also includes the author
type ResponseCollaborationTodo struct {
	Id            int64                  `json:"id"`
	Title         string                 `json:"title"`
	Body          string                 `json:"body"`
	UpdatedAt     time.Time              `json:"updated_at"`
	CreatedAt     time.Time              `json:"created_at"`
	Collaborators []ResponseCollaborator `json:"collaborators"`
	Author        ResponseCollaborator   `json:"author"`
}

type ResponseCollaborations struct {
	Data []ResponseCollaborationTodo
}
