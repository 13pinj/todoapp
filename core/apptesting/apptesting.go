package apptesting

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
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
	r.GET("/test/url", fn)
	r.POST("/test/url", fn)

	addr := fmt.Sprintf(":%v", lastPort)
	urlString := fmt.Sprintf("http://localhost:%v/test/url", lastPort)
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

// Client - тестовый клиент, основанный на http.Client, но автоматически
// сохраняющий все куки, полученные с сервера.
type Client struct {
	*http.Client
}

// NewClient инициализует клиент с пустым CookieJar
func NewClient() *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{&http.Client{Jar: jar}}
}

// Get ведет себя как http.Client.Get, но в случае успешного запроса
// сохраняет все полученные куки.
func (c *Client) Get(rawURL string) (resp *http.Response, err error) {
	resp, err = c.Client.Get(rawURL)
	if err != nil {
		return
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	c.Client.Jar.SetCookies(parsed, resp.Cookies())
	return
}

// PostForm ведет себя как http.Client.PostForm, но в случае успешного запроса
// сохраняет все полученные куки.
func (c *Client) PostForm(rawURL string, data url.Values) (resp *http.Response, err error) {
	resp, err = c.Client.PostForm(rawURL, data)
	if err != nil {
		return
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	c.Client.Jar.SetCookies(parsed, resp.Cookies())
	return
}

// ClearCookie стирает все куки, хранящиеся в клиенте.
func (c *Client) ClearCookie() {
	c.Client.Jar, _ = cookiejar.New(nil)
}
