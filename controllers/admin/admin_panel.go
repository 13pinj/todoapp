package admin

import (
	"strconv"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers"
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
