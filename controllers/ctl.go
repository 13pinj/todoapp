package ctl

import (
	"net/http"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/models/user"
)

// RenderHTML отвечает на запрос кодом 200 и выполненным шаблоном.
func RenderHTML(c *gin.Context, template string, data gin.H) {
	u, ok := user.FromContext(c)
	if data == nil {
		data = gin.H{}
	}
	data["LoggedIn"] = ok
	if ok {
		data["CurrentUser"] = u
	}
	c.HTML(http.StatusOK, template, data)
}

// RenderJSON отвечает на запрос кодом 200 и данными data, закодированными в JSON.
func RenderJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

const errTmpl = "error.tmpl"

// Render403 отвечает кодом 403 Forbidden.
func Render403(c *gin.Context) {
	c.HTML(http.StatusForbidden, errTmpl, gin.H{
		"Code": 403,
		"Msg":  "Доступ запрещен :(",
	})
}

// Render404 отвечает кодом 404 Not Found.
func Render404(c *gin.Context) {
	c.HTML(http.StatusNotFound, errTmpl, gin.H{
		"Code": 404,
		"Msg":  "У нас нет такой странички :(",
	})
}

// Render500 отвечает кодом 500 Internal Server Error.
func Render500(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, errTmpl, gin.H{
		"Code": 500,
		"Msg":  "Внутренняя ошибка!",
	})
}

// Redirect выполняет перенаправление.
func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}
