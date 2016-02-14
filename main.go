package main

import (
	"os"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers/page"
	"github.com/13pinj/todoapp/controllers/todos"
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/13pinj/todoapp/core/log"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.Writer()))
	r.Use(gin.RecoveryWithWriter(log.Writer()))

	r.LoadHTMLGlob("templates/*")
	r.Static("/s", "public")
	r.GET("/", page.Home)

	r.POST("/login", users.Login)
	r.GET("/register", users.RegistrationForm)
	r.POST("/register", users.Register)
	r.POST("/logout", users.Logout)

	r.POST("/list-create", todos.CreateList)
	r.GET("/list/:id", todos.ShowList)
	r.POST("/list/:id/update", todos.UpdateList)
	r.POST("/list/:id/destroy", todos.DestroyList)
	r.POST("/list/:id/add", todos.CreateTask)
	r.POST("/task/:id/update", todos.UpdateTask)
	r.POST("/task/:id/destroy", todos.DestroyTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
