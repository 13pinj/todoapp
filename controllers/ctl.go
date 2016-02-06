package ctl

import "github.com/gin-gonic/gin"

// Отвечает на запрос кодом 200 и выполненным шаблоном.
func RenderHTML(c *gin.Context, template string, data interface{}) {

}

// Отвечает на запрос кодом 200 и данными data, закодированными в JSON.
func RenderJSON(c *gin.Context, data interface{}) {

}

// Отвечает кодом 403 Forbidden.
func Render403(c *gin.Context) {

}

// Отвечает кодом 404 Not Found.
func Render404(c *gin.Context) {

}

// Отвечает кодом 500 Internal Server Error.
func Render500(c *gin.Context) {

}
