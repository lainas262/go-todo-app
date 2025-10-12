package main

import (
	"context"
	"example/todolist/handler"
	"example/todolist/repository"
	"example/todolist/router"
	"example/todolist/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coder/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool
var WS *websocket.Conn

// type APIError struct {
// 	Status  int    `json:"-"`
// 	Message string `json:"message"`
// }

// func (e *APIError) Error() string {
// 	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Message)
// }

// func WriteError(w http.ResponseWriter, err *APIError) {
// 	w.WriteHeader(err.Status)
// 	w.Header().Set("Content-Type", "application/json")
// 	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
// 		log.Printf("Failed to write error response: %v", encodeErr)
// 	}
// }

// var toDoItems = []model.Todo{}

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.Atoi(vars["userId"])
// 	var user model.User

// 	rows, err := DB.Query(context.Background(), "select * from users where id=$1", id)

// 	if err != nil {
// 		log.Printf("Error quering user id '%v' : %v", id, err)
// 	}

// 	user, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			WriteError(w, &APIError{404, "Not Found"})
// 			return
// 		}
// 		WriteError(w, &APIError{500, "Server Error"})
// 		log.Printf("Error collecting user id '%v' : %v", id, err)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

// func getUsers(w http.ResponseWriter, r *http.Request) {
// 	rows, err := DB.Query(context.Background(), "select * from users")
// 	if err != nil {
// 		log.Printf("Error quering users: %v", err)
// 		WriteError(w, &APIError{500, "Server Error"})
// 		return
// 	}
// 	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
// 	if err != nil {
// 		log.Printf("Error collecting rows: %v", err)
// 		WriteError(w, &APIError{500, "Server Error"})
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(users)
// }

// func createUser(w http.ResponseWriter, r *http.Request) {
// 	var newUserData model.User
// 	var newId int64

// 	err := json.NewDecoder(r.Body).Decode(&newUserData)
// 	if err != nil {
// 		WriteError(w, &APIError{500, "Invalid JSON"})
// 		return
// 	}

// 	err = DB.QueryRow(context.Background(), `
// 		INSERT INTO users (first_name, last_name, email, username)
// 		VALUES ($1, $2, $3, $4)
// 		RETURNING id`, newUserData.FirstName, newUserData.LastName, newUserData.Email, newUserData.UserName).Scan(&newId)

// 	if err != nil {
// 		log.Printf("Error inserting user: %v", err)
// 		WriteError(w, &APIError{500, "Insert Failed"})
// 		return
// 	}
// 	// if int(cmdTag.RowsAffected()) != 1 {
// 	// 	log.Printf("Warning: Expected 1 row affected, got %d", cmdTag.RowsAffected())
// 	// }
// 	newUserData.Id = newId
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newUserData)
// }

// func getTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	id, _ := strconv.Atoi(vars["id"])

// 	for _, v := range toDoItems {
// 		if v.Id == int64(id) {
// 			w.Header().Set("Content-Type", "application/json")
// 			encoder := json.NewEncoder(w)
// 			encoder.Encode(v)
// 			return
// 		}
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusNotFound)
// 	json.NewEncoder(w).Encode("Not Found")
// }

// func getTodos(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(toDoItems)
// }

// func deleteTodoCollaborator(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	// userId, _ := strconv.Atoi(vars["userId"])
// 	todoId, _ := strconv.Atoi(vars["todoId"])
// 	collaboratorId, _ := strconv.Atoi(vars["collaboratorId"])

// 	cmdTag, err := DB.Exec(context.Background(), `
// 		delete from todo_collaborators where todo_id=$1 and user_id=$2
// 	`, todoId, collaboratorId)

// 	if err != nil {
// 		log.Printf("Error deleting collaborator: %v", err)
// 		WriteError(w, &APIError{500, "Delete Failed"})
// 		return
// 	}
// 	if cmdTag.RowsAffected() == 0 {
// 		WriteError(w, &APIError{404, "Not found"})
// 		return
// 	}

// 	if cmdTag.RowsAffected() > 1 {
// 		log.Printf("Warning: Expected 1 row affected, got %d", cmdTag.RowsAffected())
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	if cmdTag.RowsAffected() == 1 {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusNoContent)
// 		return
// 	}
// }

// func getUserCollaborationTodos(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])

// 	rows, err := DB.Query(context.Background(), `
// 	SELECT t.id, t.title, t.body, t.created_at, t.updated_at, t.author_id from todo_collaborators tc
// 	INNER JOIN todos t
// 	ON tc.todo_id = t.id
// 	WHERE tc.user_id = $1
// 	`, userId)

// 	if err != nil {
// 		log.Printf("Error quering user collaborations: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	collaborationTodos, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])

// 	if err != nil {
// 		log.Printf("Error collecting rows: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(collaborationTodos)
// }

// func updateCollaborationTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	todoId, _ := strconv.Atoi(vars["todoId"])
// 	var todoItem model.Todo

// 	err := json.NewDecoder(r.Body).Decode(&todoItem)

// 	if err != nil {
// 		WriteError(w, &APIError{500, "Invalid JSON"})
// 	}

// 	todoItem.Id = int64(todoId)
// 	todoItem.UpdatedAt = time.Now()

// 	err = DB.QueryRow(context.Background(), `
// 		update todos
// 		set title=$1, body=$2, updated_at=$3
// 		where id=$4
// 		returning created_at, author_id
// 	`, todoItem.Title, todoItem.Body, todoItem.UpdatedAt, todoItem.Id).Scan(&todoItem.CreatedAt, &todoItem.Author)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			WriteError(w, &APIError{404, "Not Found"})
// 			return
// 		}

// 		log.Printf("Error updating row: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(todoItem)
// }

func ConnectDB() {

	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), // Should be the Docker service name 'db'
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database connection check failed after startup: %v", err)
	}
	DB = pool
	log.Println("Database connection pool successfully established")
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file. Falling back to system environment variables: %v", err)
	}

	//setup db connection pool
	ConnectDB()
	defer DB.Close()
	//repositories
	userRepo := repository.CreateUserRepository(DB)
	todoRepo := repository.CreateTodoRepository(DB)
	//services
	userService := service.CreateUserService(userRepo)
	todoService := service.CreateTodoService(todoRepo)

	//handlers
	userHandler := handler.CreateUserHandler(userService)
	todoHandler := handler.CreateTodoHandler(todoService)

	handler := router.SetupRouter(&router.Handlers{UserHandler: userHandler, TodoHandler: todoHandler})
	http.ListenAndServe(":8080", handler)
}
