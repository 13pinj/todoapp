package ctl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RenderHTML отвечает на запрос кодом 200 и выполненным шаблоном.
func RenderHTML(c *gin.Context, template string, data interface{}) {
	c.HTML(http.StatusOK, template, data)
}

// RenderJSON отвечает на запрос кодом 200 и данными data, закодированными в JSON.
func RenderJSON(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// Render403 отвечает кодом 403 Forbidden.
func Render403(c *gin.Context) {
	c.String(http.StatusForbidden, "403 Forbidden")
}

// Render404 отвечает кодом 404 Not Found.
func Render404(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}

// Render500 отвечает кодом 500 Internal Server Error.
func Render500(c *gin.Context) {
	c.String(http.StatusInternalServerError, "500 Internal Server Error")
}

// Redirect выполняет перенаправление.
func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}
