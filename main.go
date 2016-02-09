package main

import (
	"github.com/13pinj/todoapp/controllers/page"
	"github.com/13pinj/todoapp/controllers/todos"
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/gin-gonic/gin"
)

func main() {
	site := gin.Default()
	site.LoadHTMLGlob("templates/*")
	site.Static("/s", "public")
	site.GET("/", page.Home)

	site.POST("/login", users.Login)
	site.GET("/register", users.RegistrationForm)
	site.POST("/register", users.Register)
	site.POST("/logout", users.Logout)

	site.POST("/list-create", todos.CreateList)
	site.GET("/list/:id", todos.ShowList)
	site.POST("/list/:id/update", todos.UpdateList)
	site.POST("/list/:id/destroy", todos.DestroyList)
	site.POST("/list/:id/add", todos.CreateTask)
	site.POST("/task/:id/update", todos.UpdateTask)
	site.POST("/task/:id/destroy", todos.DestroyTask)

	site.Run(":8080")
}
