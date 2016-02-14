package users

import (
	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/controllers"
	"github.com/13pinj/todoapp/models/user"
)

// RegistrationForm отсылает форму регистрации.
// Пользователей, уже выполнивших вход, она перенаправляет на главную.
// GET /register
func RegistrationForm(c *gin.Context) {
	ctl.RenderHTML(c, "register.tmpl", nil)
}

// Register регистрирует нового пользователя, используя параметры POST-запроса
// name и password, и выполняет вход под его именем.
// Пользователей, уже выполнивших вход, она перенаправляет на главную,
// не перезаписывая сессию.
// POST /register
func Register(c *gin.Context) {
	_, ok := user.FromContext(c)
	if ok {
		ctl.Redirect(c, "/")
		return
	}
	name := c.PostForm("name")
	pas := c.PostForm("password")
	us, errs := user.Register(name, pas)
	if errs != nil {
		ctl.RenderHTML(c, "register.tmpl", gin.H{"Errors": errs})
		return
	}
	us.AutoLogin(c)
	ctl.Redirect(c, "/")
}

// LoginForm отсылает форму входа.
// Пользователей, уже выполнивших вход, она перенаправляет на главную.
// GET /
func LoginForm(c *gin.Context) {
	ctl.RenderHTML(c, "login.tmpl", nil)
}

// Login выполняет вход, используя параметры POST-запроса
// name и password.
// Пользователей, уже выполнивших вход, она перенаправляет на главную,
// не перезаписывая сессию.
// POST /login
func Login(c *gin.Context) {
	_, ok := user.FromContext(c)
	if ok {
		ctl.Redirect(c, "/")
		return
	}
	name := c.PostForm("name")
	pas := c.PostForm("password")
	_, ok = user.Login(c, name, pas)
	if !ok {
		ctl.RenderHTML(c, "login.tmpl", gin.H{"Errors": []string{"Ошибка авторизации"}})
		return
	}
	ctl.Redirect(c, "/")
}

// Logout выполняет выход и перенаправляет на главную.
// POST /logout
func Logout(c *gin.Context) {
	user.Logout(c)
	ctl.Redirect(c, "/")
}

// Destroy удаляет текущего пользователя и перенаправляет на главную.
// POST /user/destroy
func Destroy(c *gin.Context) {
	us, ok := user.FromContext(c)
	if ok {
		us.Destroy()
		ctl.Redirect(c, "/")
	}
	ctl.Redirect(c, "/")
}
