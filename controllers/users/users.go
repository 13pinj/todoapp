package users

import "github.com/gin-gonic/gin"

// RegistrationForm отсылает форму регистрации.
// Пользователей, уже выполнивших вход, она перенаправляет на главную.
// GET /register
func RegistrationForm(c *gin.Context) {

}

// Register регистрирует нового пользователя, используя параметры POST-запроса
// name и password, и выполняет вход под его именем.
// Пользователей, уже выполнивших вход, она перенаправляет на главную,
// не перезаписывая сессию.
// POST /register
func Register(c *gin.Context) {

}

// LoginForm отсылает форму входа.
// Пользователей, уже выполнивших вход, она перенаправляет на главную.
// GET /login
func LoginForm(c *gin.Context) {

}

// Login выполняет вход, используя параметры POST-запроса
// name и password.
// Пользователей, уже выполнивших вход, она перенаправляет на главную,
// не перезаписывая сессию.
// POST /login
func Login(c *gin.Context) {

}

// Logout выполняет выход и перенаправляет на главную.
// POST /logout
func Logout(c *gin.Context) {

}

// Destroy удаляет текущего пользователя и перенаправляет на главную.
// POST /user/destroy
func Destroy(c *gin.Context) {

}
