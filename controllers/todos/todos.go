package todos

import (
	"net/http"
	"strconv"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers"
	"github.com/13pinj/todoapp/models/todo"
	"github.com/13pinj/todoapp/models/todolist"
	"github.com/13pinj/todoapp/models/user"
)

// Index выводит страницу со всеми списками дел текущего пользователя.
// Если пользователь незалогинен, перенаправляет на главную.
// GET /
func Index(c *gin.Context) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Redirect(c, "/")
		return
	}
	u.LoadLists()
	ctl.RenderHTML(c, "todos_index.tmpl", gin.H{
		"Lists": u.Lists,
	})
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
		u.LoadLists()
		ctl.RenderHTML(c, "todos_index.tmpl", gin.H{
			"Lists":      u.Lists,
			"AlertError": err.Error(),
		})
		return
	}
	ctl.Redirect(c, l.Path())
}

// ShowList выводит страницу списка дел, на которой отображается его заголовок
// и содержание.
// GET /list/:id
func ShowList(c *gin.Context) {
	l, ok := getlist(c)
	if !ok {
		return
	}
	u, _ := user.FromContext(c)
	u.LoadLists()
	l.LoadTodos()
	ctl.RenderHTML(c, "todos_show.tmpl", gin.H{
		"List":  l,
		"Lists": u.Lists,
	})
}

// UpdateList изменяет заголовок списка на тот, который был получен
// POST-параметре title и перенаправляет на страницу этого списка.
// POST /list/:id/update
func UpdateList(c *gin.Context) {
	l, ok := getlist(c)
	if !ok {
		return
	}
	l.Title = c.PostForm("title")
	err := l.Save()
	if err != nil {
		ctl.RenderJSON(c, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	ctl.RenderJSON(c, gin.H{
		"status": "success",
	})
}

// DestroyList стирает список из базы и перенаправляет на главную.
// POST /list/:id/destroy
func DestroyList(c *gin.Context) {
	l, ok := getlist(c)
	if !ok {
		return
	}
	l.Destroy()
	c.Status(http.StatusOK)
}

// CreateTask создает новое задание в списке с текстом
// из POST-параметра label и перенаправляет на страницу списка.
// POST /list/:id/add
func CreateTask(c *gin.Context) {
	l, ok := getlist(c)
	if !ok {
		return
	}
	err := l.Add(c.PostForm("label"))
	if err != nil {
		u, _ := user.FromContext(c)
		u.LoadLists()
		l.LoadTodos()
		ctl.RenderHTML(c, "todos_show.tmpl", gin.H{
			"List":       l,
			"AlertError": err.Error(),
			"Lists":      u.Lists,
		})
		return
	}
	ctl.Redirect(c, l.Path())
}

// UpdateTask изменяет поля задания используя POST-параметры done и label.
// Поля, для который не заданы значения в параметрах запроса, должны остаться
// неизменными.
// После выполнения запроса UpdateTask перенаправляет клиент на страницу списка.
// POST /task/:id/update
func UpdateTask(c *gin.Context) {
	td, _, ok := gettask(c)
	if !ok {
		return
	}
	label, ok := c.GetPostForm("label")
	if ok {
		td.Label = label
	}
	done, ok := c.GetPostForm("done")
	if ok {
		td.Done = (done != "0")
	}
	err := td.Save()
	if err != nil {
		ctl.RenderJSON(c, gin.H{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	ctl.RenderJSON(c, gin.H{
		"status": "success",
	})
}

// DestroyTask стирает задание из списка и перенаправляет на страницу списка.
// POST /task/:id/destroy
func DestroyTask(c *gin.Context) {
	td, _, ok := gettask(c)
	if !ok {
		return
	}
	td.Destroy()
	c.Status(http.StatusOK)
}

func getlist(c *gin.Context) (*todolist.TodoList, bool) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Render403(c)
		return nil, false
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctl.Render404(c)
		return nil, false
	}
	l, ok := todolist.Find(uint(id))
	if !ok {
		ctl.Render404(c)
		return nil, false
	}
	if l.UserID != u.ID {
		ctl.Render403(c)
		return nil, false
	}
	return l, true
}

func gettask(c *gin.Context) (*todo.Todo, *todolist.TodoList, bool) {
	u, ok := user.FromContext(c)
	if !ok {
		ctl.Render403(c)
		return nil, nil, false
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctl.Render404(c)
		return nil, nil, false
	}
	td, ok := todo.Find(uint(id))
	if !ok {
		ctl.Render404(c)
		return nil, nil, false
	}
	l, ok := todolist.Find(td.TodoListID)
	if !ok {
		ctl.Render500(c)
		return nil, nil, false
	}
	if u.ID != l.UserID {
		ctl.Render403(c)
		return nil, nil, false
	}
	return td, l, true
}
