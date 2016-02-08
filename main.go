package main

import (
	"github.com/13pinj/todoapp/controllers/page"
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/gin-gonic/gin"
)

func main() {
	site := gin.Default()
	site.LoadHTMLGlob("templates/*")
	site.GET("/", page.Home)
	site.GET("/login", users.LoginForm)
	site.POST("/login", users.Login)
	site.GET("/register", users.RegistrationForm)
	site.POST("/register", users.Register)
	site.POST("/logout", users.Logout)
	site.Run(":8080")
}
