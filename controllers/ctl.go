package ctl

import "github.com/gin-gonic/gin"

// RenderHTML отвечает на запрос кодом 200 и выполненным шаблоном.
func RenderHTML(c *gin.Context, template string, data interface{}) {

}

// RenderJSON отвечает на запрос кодом 200 и данными data, закодированными в JSON.
func RenderJSON(c *gin.Context, data interface{}) {

}

// Render403 отвечает кодом 403 Forbidden.
func Render403(c *gin.Context) {

}

// Render404 отвечает кодом 404 Not Found.
func Render404(c *gin.Context) {

}

// Render500 отвечает кодом 500 Internal Server Error.
func Render500(c *gin.Context) {

}
