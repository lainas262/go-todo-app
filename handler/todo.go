package handler

import (
	"context"
	"encoding/json"
	"example/todolist/api"
	"example/todolist/middleware"
	"example/todolist/response"
	"example/todolist/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	service *service.TodoService
}

func CreateTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) GetUserTodos(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	log.Printf("GetUserTodos - Request received: ")

	todos, err := h.service.GetUserTodos(context.Background(), userId)

	if err != nil {
		log.Printf("GetUserTodos - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GetUserTodos - Success: user id=%v", userId)

	response.WriteJSON(w, http.StatusOK, todos)

}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {

	// Get user ID from context (set by JWT middleware)
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req api.RequestCreateTodo

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateTodo - JSON decode error: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateTodo - Request received: ")

	todo, err := h.service.CreateTodo(context.Background(), userId, req)

	if err != nil {
		log.Printf("CreateTodo - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateTodo - Success: todo created with id=%v", todo.Id)

	response.WriteJSON(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {

	// Get user ID from context (set by JWT middleware)
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	vars := mux.Vars(r)
	todoId, err := strconv.Atoi(vars["todoId"])

	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid todoId parameter provided")
		return
	}

	//create and update have same request payload but different path and method on api
	var req api.RequestCreateTodo

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateTodo - JSON decode error: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("UpdateTodo - Request received: ")

	todo, err := h.service.UpdateTodo(context.Background(), userId, int64(todoId), req)

	if err != nil {
		log.Printf("UpdateTodo - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("UpdateTodo - Success: todo updated with id=%v", todo.Id)

	response.WriteJSON(w, http.StatusOK, todo)
}
