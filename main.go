package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SomeFirstPage(c *gin.Context) {
	c.String(http.StatusOK, "Ну здраствуй,\n")
	c.String(http.StatusOK, "Ты на заглавной странице сайта,\n")
	c.String(http.StatusOK, "На котором пока что совсем ничего нету")
}

func main() {
	site := gin.Default()
	site.LoadHTMLGlob("templates/*")
	site.GET("/", SomeFirstPage)
	site.Run(":8080")
}
