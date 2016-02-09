package page

import (
	"github.com/13pinj/todoapp/controllers/todos"
	"github.com/13pinj/todoapp/controllers/users"
	"github.com/13pinj/todoapp/models/user"
	"github.com/gin-gonic/gin"
)

// Home выводит нужную главную страницу в зависимости от статуса авторизации
// пользователя. Авторизованные пользователи попадают на страницу
// со своими списками дел, остальные на страницу входа.
func Home(c *gin.Context) {
	if _, ok := user.FromContext(c); ok {
		todos.Index(c)
	} else {
		users.LoginForm(c)
	}
}
