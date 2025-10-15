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

	wsHandler "example/todolist/ws/handler"
	wsService "example/todolist/ws/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

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

	//ws
	wsService := wsService.CreateWebsocketService()
	go wsService.StartHub()
	wsHandler := wsHandler.CreateWebsocketHandler(wsService)

	handler := router.SetupRouter(&router.Handlers{UserHandler: userHandler, TodoHandler: todoHandler, WsHandler: wsHandler})
	http.ListenAndServe(":8080", handler)
}
