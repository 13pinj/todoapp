package page

import (
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	users.LoginForm(c)
}
