package handler

import (
	"context"
	"encoding/json"
	"example/todolist/api"
	"example/todolist/response"
	"example/todolist/service"
	"log"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func CreateUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req api.RequestCreateUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateUser - JSON decode error: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
	}

	log.Printf("CreateUser - Request received: username=%s, email=%s", req.Username, req.Email)

	user, err := h.service.CreateUser(context.Background(), req)

	if err != nil {
		log.Printf("CreateUser - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateUser - Success: user created with username=%s", user.Username)

	response.WriteJSON(w, http.StatusOK, user)
}
