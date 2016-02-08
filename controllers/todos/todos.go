package todos

import (
	"net/http"
	"strconv"

	"github.com/13pinj/todoapp/controllers"
	"github.com/13pinj/todoapp/models/todolist"
	"github.com/13pinj/todoapp/models/user"
	"github.com/gin-gonic/gin"
)

// Index выводит страницу со всеми списками дел текущего пользователя.
// Если пользователь незалогинен, перенаправляет на главную.
// GET /
func Index(c *gin.Context) {

}

// CreateList создает новый список дел с заголовком из POST-параметра title
// и перенаправляет на страницу этого списка.
// POST /list-create
func CreateList(c *gin.Context) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Render403(c)
		return
	}
	l := todolist.New(c.PostForm("title"))
	l.UserID = u.ID
	err := l.Save()
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	ctl.Redirect(c, l.Path())
}

// ShowList выводит страницу списка дел, на которой отображается его заголовок
// и содержание.
// GET /list/:id
func ShowList(c *gin.Context) {

}

// UpdateList изменяет заголовок списка на тот, который был получен
// POST-параметре title и перенаправляет на страницу этого списка.
// POST /list/:id/update
func UpdateList(c *gin.Context) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Render403(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctl.Render404(c)
		return
	}
	l, ok := todolist.Find(uint(id))
	if !ok {
		ctl.Render404(c)
		return
	}
	if l.UserID != u.ID {
		ctl.Render403(c)
		return
	}
	l.Title = c.PostForm("title")
	err = l.Save()
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	ctl.Redirect(c, l.Path())
}

// DestroyList стирает список из базы и перенаправляет на главную.
// POST /list/:id/destroy
func DestroyList(c *gin.Context) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Render403(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctl.Render404(c)
		return
	}
	l, ok := todolist.Find(uint(id))
	if !ok {
		ctl.Render404(c)
		return
	}
	if l.UserID != u.ID {
		ctl.Render403(c)
		return
	}
	l.Destroy()
	ctl.Redirect(c, "/")
}

// CreateTask создает новое задание в списке с текстом
// из POST-параметра label и перенаправляет на страницу списка.
// POST /list/:id/task-create
func CreateTask(c *gin.Context) {

}

// UpdateTask изменяет поля задания используя POST-параметры done и label.
// Поля, для который не заданы значения в параметрах запроса, должны остаться
// неизменными.
// После выполнения запроса UpdateTask перенаправляет клиент на страницу списка.
// POST /list/:id/task-update/:task-id
func UpdateTask(c *gin.Context) {

}

// DestroyTask стирает задание из списка и перенаправляет на страницу списка.
// POST /list/:id/task-destroy/:task-id
func DestroyTask(c *gin.Context) {

}
