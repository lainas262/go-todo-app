package router

import (
	"example/todolist/handler"
	"example/todolist/middleware"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Handlers struct {
	UserHandler *handler.UserHandler
	TodoHandler *handler.TodoHandler
	WsHandler   *handler.WebsocketHandler
}

func SetupRouter(h *Handlers) http.Handler {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/login", h.UserHandler.Login).Methods("POST")

	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("", h.UserHandler.CreateUser).Methods("POST")
	usersSelfRouter := usersRouter.PathPrefix("/self").Subrouter()
	usersSelfRouter.Use(middleware.JWTAuth)
	usersSelfRouter.HandleFunc("", h.UserHandler.GetSelf).Methods("GET")

	todosRouter := apiRouter.PathPrefix("/todos").Subrouter()
	todosRouter.Use(middleware.JWTAuth)
	todosRouter.HandleFunc("", h.TodoHandler.GetUserTodos).Methods("GET")
	todosRouter.HandleFunc("", h.TodoHandler.CreateTodo).Methods("POST")
	todosRouter.HandleFunc("/{todoId}", h.TodoHandler.UpdateTodo).Methods("PATCH")
	todosRouter.HandleFunc("/{todoId}", h.TodoHandler.DeleteTodo).Methods("DELETE")
	todosRouter.HandleFunc("/{todoId}/collaborators", h.TodoHandler.UpdateCollaborators).Methods("PATCH")
	todosRouter.HandleFunc("/{todoId}/collaborators", h.TodoHandler.DeleteCollaborator).Methods("DELETE")

	collaborationsRouter := apiRouter.PathPrefix("/collaborations").Subrouter()
	collaborationsRouter.Use(middleware.JWTAuth)
	collaborationsRouter.HandleFunc("", h.TodoHandler.GetUserCollaborationTodos).Methods("GET")
	collaborationsRouter.HandleFunc("/{todoId}", h.TodoHandler.UpdateCollaborationTodo).Methods("PATCH")

	wsRouter := r.PathPrefix("/subscribe").Subrouter()
	wsRouter.Use(middleware.WbsocketTokenAuth)
	wsRouter.HandleFunc("", h.WsHandler.Subscribe)

	// Add CORS headers
	// This allows your React app at localhost:3000 to make requests to your Go API
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r)

	return corsHandler
}
