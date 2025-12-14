package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"privateCabin/handler"
	"privateCabin/repository"
	"privateCabin/service"
	"syscall"
	"time"

	_ "privateCabin/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Startup_PrivateCabin API
// @version 1.0
// @description API для проверки регистрации и входа

// @host localhost:8080
// @BasePath /
func main() {

	//OPEN BD
	_ = godotenv.Load()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		//connStr = "postgres://user:pass@localhost:5433/mydb?sslmode=disable"
		log.Fatal("DB URL NOT SET")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT IS NOT SET")
	}
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatalf("Ошибка подключения к базе: %v", err)
	}

	// Проверяем, что соединение живое
	if err := db.Ping(); err != nil {
		log.Fatalf("БД недоступна: %v", err)
	}

	log.Println("Подключено к базе")

	//CREATE ADAPTERP
	repo := repository.NewUserRepository(db)
	srv := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(srv)
	//////////////////

	//NEW ROUTER
	r := chi.NewRouter()
	//MIDDLEWARE
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:" + port, "http://127.0.0.1:" + port},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	//swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	//GROUP ROUTE
	r.Route("/", func(r chi.Router) {

		r.Post("/login", userHandler.GetUser)
		r.Post("/register", userHandler.CreateUser)
	})

	// Graceful shutdown должен быть ПЕРЕД http.ListenAndServe
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	server := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		log.Println("server started on localhost:" + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидаем сигнала
	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down...\n", sig)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v\n", err)
	}
	log.Println("Graceful shutdown complete")
}
