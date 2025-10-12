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
		return
	}

	log.Printf("CreateUser - Request received: username=%s, email=%s", req.Username, req.Email)

	user, err := h.service.CreateUser(context.Background(), req)

	if err != nil {
		log.Printf("CreateUser - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateUser - Success: user created with username=%s", user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  user.AccessToken,
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	})

	response.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req api.RequestLoginUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("LoginUser - JSON decode error: %v", err)
		response.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("LoginUser - Request received: email=%s", req.Email)

	user, err := h.service.Login(context.Background(), req)

	if err != nil {
		log.Printf("LoginUser - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("LoginUser - Success: user logged in with username=%s", user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		Value:  user.AccessToken,
		MaxAge: 60 * 60 * 24,
		Path:   "/",
	})

	response.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetSelf(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(middleware.UserIdKey).(int64)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	log.Printf("GetSelf - Request received: ")

	user, err := h.service.GetSelf(context.Background(), userId)

	if err != nil {
		log.Printf("GetSelf - Service error: %v", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("GetSelf - Success: user fetched with username=%s", user.Username)

	response.WriteJSON(w, http.StatusOK, user)
}
