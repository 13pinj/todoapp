package apptesting

import (
	"fmt"
	"net/url"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Server - структура тестового сервера.
// Это сервер gin, работающий в сторонней горутине и принимающий
// подключения только по одному пути и работающий только с одной HandlerFunc.
// Используется для тестирования серверных компонентов путем имитации входящих
// подключений.
type Server struct {
	// URL, по которому принимает подключения сервер.
	URL *url.URL
}

var lastPort = 8080

// NewServer запускает тестовый сервер на случайном порте, который направляет
// входящие подключения на функцию fn.
func NewServer(fn gin.HandlerFunc) *Server {
	// Чтобы тестирования не началось случайно до того, как запустится сервер Gin.
	runtime.GOMAXPROCS(1)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", fn)

	addr := fmt.Sprintf(":%v", lastPort)
	urlString := fmt.Sprintf("http://localhost:%v/", lastPort)
	lastPort++

	// Сервер запускается в отдельной горутине.
	go r.Run(addr)
	// Передаем очередь выполнения горутине сервера.
	runtime.Gosched()
	// Контроль текущей горутине вернется, как только сервер заблокируется
	// в ожидании подключений - то есть когда он будет полностью инициализирован.
	// Довольно костыльное решение, но gin не предоставляет
	// нормального функционала для тестирования.

	serverURL, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}
	return &Server{URL: serverURL}
}
