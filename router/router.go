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

	// todosRouter.HandleFunc("", getTodos).Methods("GET")

	// usersRouter.HandleFunc("", getUsers).Methods("GET")

	// usersRouter.HandleFunc("/{userId}", getUser).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos", getUserTodos).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", getUserTodo).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", deleteUserTodo).Methods("DELETE")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", updateUserTodo).Methods("PATCH")

	// usersRouter.HandleFunc("/{userId}/todos/{todoId}/collaborators", updateUserTodoCollaborators).Methods("PATCH")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}/collaborators/{collaboratorId}", deleteTodoCollaborator).Methods("DELETE")

	// usersRouter.HandleFunc("/{userId}/collaborations", getUserCollaborationTodos).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/collaborations/{todoId}", updateCollaborationTodo).Methods("PATCH")

	// todosRouter.HandleFunc("/{id}", getTodo).Methods("GET")

	// r.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
	// 	wsConn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
	// 		OriginPatterns: []string{"http://localhost:3000"},
	// 	})

	// 	WS = wsConn

	// 	if err != nil {
	// 		log.Printf("Websocket accept err: %v", err)
	// 		return
	// 	}
	// 	defer WS.CloseNow()

	// 	for {
	// 		var v any
	// 		err = wsjson.Read(context.Background(), WS, &v)
	// 		if err != nil {
	// 			log.Printf("Error reading message: %v", err)
	// 			break
	// 		}
	// 		log.Println(v)

	// 	}

	// })

	// Add CORS headers
	// This allows your React app at localhost:3000 to make requests to your Go API
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)(r)

	return corsHandler
}
