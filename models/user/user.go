package user

import (
	"github.com/13pinj/todoapp/models/todolist"
	"github.com/gin-gonic/gin"
)

// User - структура модели пользователя
type User struct {
	ID      uint
	Name    string
	PwdHash string
	Lists   []*todolist.TodoList
}

// imUser - внутреннее представление модели.
type imUser struct {
	ID      uint
	Name    string
	PwdHash string
}

// Хранилище моделей в памяти.
var imStorage = make(map[uint]*imUser)

// Формутирует внутреннее представление модели во внешнее.
func (u *imUser) format() *User {
	return nil
}

// Преобразует внешнее представление во внутреннее.
func (u *User) zip() *imUser {
	return nil
}

// validate проверяет корректность всех полей модели.
func (u *User) validate() error {
	return nil
}

// Register добавляет нового пользователя в базу, и возвращает его структуру,
// если введенные поля корректны. В противном случае Register возвращает ошибку.
func Register(name string, password string) (*User, error) {
	return nil, nil
}

// Login выполняет авторизацию пользователей.
// Если введенные имя и пароль действительны, Login запишет факт авторизации
// в сессию пользователя и вернет первым значением структуру пользователя,
// а вторым true. В противном случае - nil и false.
// Login перезапишет старые данные об авторизации, если таковые имеются.
func Login(c *gin.Context, name string, password string) (*User, bool) {
	return nil, false
}

// Logout стирает данные об авторизации из сессии пользователя.
func Logout(c *gin.Context) {

}

// FromContext получает данные об авторизации из сессии пользователя.
// Если пользователь авторизован, FromContext вернет структуру и true.
// Иначе nil и false.
func FromContext(c *gin.Context) (*User, bool) {
	return nil, false
}

// AutoLogin запишет факт авторизации в сессию пользователя.
// Он перезапишет старые данные об авторизации, если таковые имеются.
func (u *User) AutoLogin(c *gin.Context) {

}

// Destroy стирает данные о пользователе из базы данных.
func (u *User) Destroy() {

}
