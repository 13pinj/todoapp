package admin

import (
	"strconv"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers"
	"github.com/13pinj/todoapp/models"
	"github.com/13pinj/todoapp/models/todo"
	"github.com/13pinj/todoapp/models/todolist"
	"github.com/13pinj/todoapp/models/user"
)

const (
	usersPerPage = 15
)

// Index возвращает главную страницу админской панели.
// GET /admin
func Index(c *gin.Context) {
	if !assertAccess(c) {
		return
	}
	page, err := strconv.Atoi(c.Query("p"))
	if err != nil {
		page = 1
	}
	pcount := user.Pages(usersPerPage)
	ctl.RenderHTML(c, "admin_index.tmpl", gin.H{
		"Users": user.FindPage(page, usersPerPage, user.ByVisitedAtDesc),
		"Pager": gin.H{
			"Cur":      page,
			"Max":      pcount,
			"PathTmpl": "/admin?p=",
		},
	})
}

// User возвращает страницу с информацией о пользователе.
// GET /admin/u/:name
func User(c *gin.Context) {
	if !assertAccess(c) {
		return
	}
	u, ok := user.Find(c.Param("name"))
	if !ok {
		ctl.Render404(c)
		return
	}
	u.LoadLists()
	for _, i := range u.Lists {
		i.LoadTodos()
	}
	// Шаблон ожидает заполненую структуру пользователя
	ctl.RenderHTML(c, "admin_user.tmpl", gin.H{
		"User": u,
	})
}

// UserUpdate обновляет информацию о пользователе.
// POST /admin/u/:name
func UserUpdate(c *gin.Context) {

}

// UserDestroy стирает пользователя из базы.
// POST /admin/u/:name/destroy
func UserDestroy(c *gin.Context) {

}

func assertAccess(c *gin.Context) bool {
	u, ok := user.FromContext(c)
	if !ok || !u.Admin() {
		ctl.Render403(c)
		return false
	}
	return true
}

// TrashStatus - подсчитывает число удаленных записей
func TrashStatus() (notdel int, del int) {
	var buff int
	models.DB.Model(&todo.Todo{}).Where("deleted_at IS NULL").Count(&notdel)
	buff = notdel
	models.DB.Model(&todolist.TodoList{}).Where("deleted_at IS NULL").Count(&notdel)
	notdel += buff
	models.DB.Model(&todo.Todo{}).Where("deleted_at IS NOT NULL").Count(&del)
	buff = del
	models.DB.Model(&todolist.TodoList{}).Where("deleted_at IS NOT NULL").Count(&del)
	del += buff
	return
}
