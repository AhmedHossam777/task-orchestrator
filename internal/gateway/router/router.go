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
			tasks.POST("/", taskHandler.CreateTaskHandler)
			tasks.GET("/", taskHandler.GetAllTasks)
			tasks.GET("/:id", taskHandler.GetOneTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
		}
	}
}
