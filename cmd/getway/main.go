package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/handler"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/middleware"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/repository"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/router"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	pool := initDB()
	defer pool.Close()
	
	taskRepo := repository.NewTaskRepository(pool)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	
	gin.SetMode(gin.ReleaseMode)
	
	engine := gin.New()
	engine.Use(
		middleware.RequestID(), middleware.Recovery(), middleware.Logger(),
		middleware.RateLimiter(20, 10),
	)
	
	router.Setup(engine, taskHandler)
	
	srv := &http.Server{
		Addr:         ":8008",
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// gracefull shutdown
	go func() {
		log.Printf("API Gateway starting on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %v, shutting down gracefully...", sig)
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	
	log.Println("Server stopped cleanly")
}

func initDB() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading the env variables: %w", err)
	}
	dsn := os.Getenv("DATABASE_URL")
	fmt.Println("dsn: ", dsn)
	if dsn == "" {
		log.Fatalf("Database url is not found")
	}
	
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
	}
	
	config.MaxConns = 25                     // Max concurrent connections
	config.MinConns = 5                      // Connections kept warm
	config.MaxConnLifetime = 5 * time.Minute // Recycle stale connections
	config.MaxConnIdleTime = 1 * time.Minute // Close idle connections
	
	// pgxpool.NewWithConfig connects immediately — no separate Ping() needed.
	// If this returns nil error, the connection is verified.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	log.Println("Connected to PostgreSQL")
	
	return pool
}
