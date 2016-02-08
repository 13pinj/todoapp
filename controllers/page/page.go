package page

import (
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/gin-gonic/gin"
)

// Home выводит нужную главную страницу в зависимости от статуса авторизации
// пользователя. Авторизованные пользователи попадают на страницу
// со своими списками дел, остальные на страницу входа.
func Home(c *gin.Context) {
	users.LoginForm(c)
}
