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

// func getUserTodos(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])

// 	rows, err := DB.Query(context.Background(), `
// 	select * from todos where author_id=$1
// 	`, userId)

// 	if err != nil {
// 		log.Printf("Error quering todos: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])

// 	if err != nil {
// 		log.Printf("Error processing rows todos: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	todoIds := make([]int64, 0, len(todos))
// 	todoIndex := make(map[int64]*model.Todo, len(todos))

// 	for i, val := range todos {
// 		todoIds = append(todoIds, val.Id)
// 		todoIndex[val.Id] = &todos[i]
// 	}

// 	collabRows, err := DB.Query(context.Background(), `
// 	select tc.todo_id, u.id, u.first_name, u.last_name, u.email, u.username
// 	FROM todo_collaborators tc
// 	INNER JOIN users u
// 	ON u.id = tc.user_id
// 	where tc.todo_id = ANY($1)
// 	`, todoIds)

// 	if err != nil {
// 		log.Printf("Error quering todo collaborators: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	defer collabRows.Close()

// 	for collabRows.Next() {
// 		var todoId int64
// 		var user model.User

// 		if err := collabRows.Scan(&todoId, &user.Id, &user.FirstName, &user.LastName, &user.Email, &user.UserName); err != nil {
// 			log.Printf("Error scanning collaborators: %v", err)
// 			WriteError(w, &APIError{500, "Internal Server Error"})
// 			return
// 		}
// 		if todo, ok := todoIndex[todoId]; ok {
// 			todo.Collaborators = append(todo.Collaborators, user)
// 		}
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(todos)
// }

// func getUserTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])
// 	todoId, _ := strconv.Atoi(vars["todoId"])

// 	var todo model.Todo

// 	err := DB.QueryRow(context.Background(), `select * from todos where author_id=$1 and id=$2
// 	`, userId, todoId).Scan(&todo.Id, &todo.Author, &todo.Body, &todo.Title, &todo.CreatedAt, &todo.UpdatedAt)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			WriteError(w, &APIError{404, "Not Found"})
// 			return
// 		}

// 		log.Printf("Error processing rows todos: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	collabRows, err := DB.Query(context.Background(), `
// 	select u.id, u.first_name, u.last_name, u.email, u.username
// 	from todo_collaborators tc
// 	INNER JOIN users u ON tc.user_id = u.id
// 	WHERE tc.todo_id = $1
// 	`, todoId)

// 	if err != nil {
// 		log.Printf("Error selecting todo collaborators: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	for collabRows.Next() {
// 		var user model.User

// 		if err := collabRows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.UserName); err != nil {
// 			log.Printf("Error parsing row todo collaborator: %v", err)
// 			WriteError(w, &APIError{500, "Internal Server Error"})
// 			return
// 		}
// 		todo.Collaborators = append(todo.Collaborators, user)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(todo)

// }

// func createUserTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])
// 	var newTodo model.Todo
// 	newTodo.Author = int64(userId)

// 	err := json.NewDecoder(r.Body).Decode(&newTodo)
// 	if err != nil {
// 		WriteError(w, &APIError{500, "Invalid JSON"})
// 		return
// 	}

// 	err = DB.QueryRow(context.Background(), `
// 		INSERT INTO todos (author_id, title, body) VALUES ($1, $2, $3)
// 		RETURNING id, created_at, updated_at`, newTodo.Author, newTodo.Title, newTodo.Body).Scan(&newTodo.Id, &newTodo.CreatedAt, &newTodo.UpdatedAt)

// 	if err != nil {
// 		log.Printf("Error inserting todo: %v", err)
// 		WriteError(w, &APIError{500, "Insert Failed"})
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newTodo)
// }

// func deleteUserTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])
// 	todoId, _ := strconv.Atoi(vars["todoId"])

// 	cmdTag, err := DB.Exec(context.Background(), `
// 		delete from todos where id=$1 and author_id=$2
// 	`, todoId, userId)

// 	if err != nil {
// 		log.Printf("Error deleting user: %v", err)
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

// func updateUserTodo(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	userId, _ := strconv.Atoi(vars["userId"])
// 	todoId, _ := strconv.Atoi(vars["todoId"])
// 	var todoItem model.Todo

// 	err := json.NewDecoder(r.Body).Decode(&todoItem)

// 	if err != nil {
// 		WriteError(w, &APIError{500, "Invalid JSON"})
// 	}

// 	todoItem.Id = int64(todoId)
// 	todoItem.Author = int64(userId)
// 	todoItem.UpdatedAt = time.Now()

// 	err = DB.QueryRow(context.Background(), `
// 		update todos
// 		set title=$1, body=$2, updated_at=$3
// 		where id=$4 and author_id=$5
// 		returning created_at
// 	`, todoItem.Title, todoItem.Body, todoItem.UpdatedAt, todoItem.Id, todoItem.Author).Scan(&todoItem.CreatedAt)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			WriteError(w, &APIError{404, "Not Found"})
// 			return
// 		}

// 		log.Printf("Error updating row: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	rows, err := DB.Query(context.Background(), `
// 	select u.id, u.username, u.email, u.first_name, u.last_name from todo_collaborators tc
// 	JOIN users u ON tc.user_id = u.id
// 	WHERE tc.todo_id=$1
// 	`, todoId)

// 	if err != nil {
// 		log.Printf("Error quering collaborators: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])

// 	if err != nil {
// 		log.Printf("Error collecting rows: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	fmt.Println(users)
// 	todoItem.Collaborators = users

// 	wsjson.Write(context.Background(), WS, todoItem)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(todoItem)
// }

// type CollaboratorsRequest struct {
// 	UserIds []int64 `json:"user_ids"`
// }

// func updateUserTodoCollaborators(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	// userId, _ := strconv.Atoi(vars["userId"])
// 	todoId, _ := strconv.Atoi(vars["todoId"])

// 	var collaboratorReq CollaboratorsRequest

// 	err := json.NewDecoder(r.Body).Decode(&collaboratorReq)

// 	if err != nil {
// 		WriteError(w, &APIError{500, "Invalid JSON"})
// 		return
// 	}

// 	if len(collaboratorReq.UserIds) == 0 {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	ctx := context.Background()
// 	tx, err := DB.Begin(ctx)

// 	if err != nil {
// 		log.Printf("Failed to begin transaction: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	defer tx.Rollback(ctx)

// 	for _, val := range collaboratorReq.UserIds {
// 		_, err := tx.Exec(context.Background(), `
// 		insert into todo_collaborators (todo_id, user_id)
// 		values ($1, $2) ON CONFLICT DO NOTHING
// 		`, todoId, val)

// 		if err != nil {
// 			log.Printf("Failed to insert collaborator %d: %v", val, err)
// 		}
// 	}
// 	if err := tx.Commit(ctx); err != nil {
// 		log.Printf("Failed to commit transaction: %v", err)
// 		WriteError(w, &APIError{500, "Internal Server Error"})
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
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
