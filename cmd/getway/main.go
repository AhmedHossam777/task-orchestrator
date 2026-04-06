package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/handler"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/repository"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/router"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/service"
	"github.com/gin-gonic/gin"
)

func main() {
	taskRepo := repository.NewInMemoryTaskRepository()
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	engine := gin.Default()
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
