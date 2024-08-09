package main

import (
	"db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func allowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func main() {
	_, err := db.OpenDb()
	if err != nil {
		panic(err)
	}
	defer db.CloseDb()

	db.CreateTables()

	r := gin.Default()

	r.Use(allowCors())

	// Project CRUD
	r.POST("/projects", db.CreateProject)
	r.GET("/projects/:id", db.GetProject)
	r.PUT("/projects/:id", db.UpdateProject)
	r.DELETE("/projects/:id", db.DeleteProject)

	// Milestone CRUD
	r.POST("/milestones", db.CreateMilestone)
	r.GET("/milestones/:id", db.GetMilestone)
	r.PUT("/milestones/:id", db.UpdateMilestone)
	r.DELETE("/milestones/:id", db.DeleteMilestone)

	// Task CRUD
	r.POST("/tasks", db.CreateTask)
	r.GET("/tasks/:id", db.GetTask)
	r.PUT("/tasks/:id", db.UpdateTask)
	r.DELETE("/tasks/:id", db.DeleteTask)

	// Subtask CRUD
	r.POST("/subtasks", db.CreateSubtask)
	r.GET("/subtasks/:id", db.GetSubtask)
	r.PUT("/subtasks/:id", db.UpdateSubtask)
	r.DELETE("/subtasks/:id", db.DeleteSubtask)

	r.GET("/projects/doing", db.GetDoingProjectsTasks)

	r.Run(":8080")
}
