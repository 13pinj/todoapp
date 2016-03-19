package main

import (
	"html/template"
	"os"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers"
	"github.com/13pinj/todoapp/controllers/admin"
	"github.com/13pinj/todoapp/controllers/page"
	"github.com/13pinj/todoapp/controllers/todos"
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/13pinj/todoapp/core/log"
	"github.com/13pinj/todoapp/templates"
)

func main() {
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.Writer()))
	r.Use(gin.RecoveryWithWriter(log.Writer()))

	tmpl, err := template.New("").Funcs(tmpl.Funcs()).ParseGlob("templates/*")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(tmpl)

	r.NoRoute(ctl.Render404)
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

	r.GET("/admin", admin.Index)
	r.GET("/admin/u/:name", admin.User)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
