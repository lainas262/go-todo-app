package router

import (
	"example/todolist/handler"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Handlers struct {
	UserHandler *handler.UserHandler
}

func SetupRouter(h *Handlers) http.Handler {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	usersRouter := apiRouter.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("", h.UserHandler.CreateUser).Methods("POST")
	// usersRouter.HandleFunc("", getUsers).Methods("GET")

	// usersRouter.HandleFunc("/{userId}", getUser).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos", getUserTodos).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos", createUserTodo).Methods("POST")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", getUserTodo).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", deleteUserTodo).Methods("DELETE")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}", updateUserTodo).Methods("PATCH")

	// usersRouter.HandleFunc("/{userId}/todos/{todoId}/collaborators", updateUserTodoCollaborators).Methods("PATCH")
	// usersRouter.HandleFunc("/{userId}/todos/{todoId}/collaborators/{collaboratorId}", deleteTodoCollaborator).Methods("DELETE")

	// usersRouter.HandleFunc("/{userId}/collaborations", getUserCollaborationTodos).Methods("GET")
	// usersRouter.HandleFunc("/{userId}/collaborations/{todoId}", updateCollaborationTodo).Methods("PATCH")

	// todosRouter := apiRouter.PathPrefix("/todos").Subrouter()
	// todosRouter.HandleFunc("", getTodos).Methods("GET")
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
