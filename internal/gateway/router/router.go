package router

import (
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/handler"
	"github.com/gin-gonic/gin"
)

func Setup(engine *gin.Engine, taskHandler *handler.TaskHandler) {
	v1 := engine.Group("/api/v1")
	{
		tasks := v1.Group("/tasks")
		{
			tasks.POST("/", taskHandler.Create)
			tasks.GET("/", taskHandler.List)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.DELETE("/:id", taskHandler.Delete)
		}
	}
	
	engine.GET(
		"/health", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"status": "ok"})
		},
	)
}
