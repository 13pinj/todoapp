package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"github.com/13pinj/todoapp/models/session"
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
var key uint

// Формутирует внутреннее представление модели во внешнее.
func (u *imUser) format() *User {
	return &User{
		ID:      u.ID,
		Name:    u.Name,
		PwdHash: u.PwdHash,
		Lists:   todolist.FindByUser(u.ID),
	}
}

// Преобразует внешнее представление во внутреннее.
func (u *User) zip() *imUser {
	return &imUser{
		ID:      u.ID,
		Name:    u.Name,
		PwdHash: u.PwdHash,
	}
}

func nameValidate(name string) bool {
	for _, v := range imStorage {
		if v.Name == name {
			return false
		}
	}
	return true
}

// Register добавляет нового пользователя в базу, и возвращает его структуру,
// если введенные поля корректны. В противном случае Register возвращает ошибку.
func Register(name string, password string) (*User, error) {
	ok := strings.Contains(name, " ")
	if ok {
		return nil, errors.New("Использование пробела в имени пользователя не допускается")
	}
	if len([]rune(name)) < 4 {
		return nil, errors.New("Использование имени пользователя короче 4 символов не допускается")
	}
	if len([]rune(password)) < 6 {
		return nil, errors.New("Использование пароля короче 4 символов не допускается")
	}
	if !nameValidate(name) {
		return nil, errors.New("Имя пользователя присутствует в базе")
	}
	record := User{}
	record.Name = name
	hash := sha1.Sum([]byte(password))
	str := fmt.Sprintf("%x", hash)
	record.PwdHash = str
	record.save()
	return &record, nil
}

func (u *User) save() {
	key++
	if u.ID == 0 {
		u.ID = key
		imStorage[u.ID] = u.zip()
	}
}

// Login выполняет авторизацию пользователей.
// Если введенные имя и пароль действительны, Login запишет факт авторизации
// в сессию пользователя и вернет первым значением структуру пользователя,
// а вторым true. В противном случае - nil и false.
// Login перезапишет старые данные об авторизации, если таковые имеются.
func Login(c *gin.Context, name string, password string) (*User, bool) {
	FoundUser := (*User)(nil)
	for _, user := range imStorage {
		if user.Name == name {
			FoundUser = user.format()
			break
		}
	}
	if FoundUser == nil {
		return nil, false
	}
	hash := sha1.Sum([]byte(password))
	str := fmt.Sprintf("%x", hash)
	if str == FoundUser.PwdHash {
		st := session.FromContext(c)
		st.SetInt("user_id", int(FoundUser.ID))
		return FoundUser, true
	}
	return nil, false
}

// Logout стирает данные об авторизации из сессии пользователя.
func Logout(c *gin.Context) {
	st := session.FromContext(c)
	st.SetInt("user_id", 0)
}

// FromContext получает данные об авторизации из сессии пользователя.
// Если пользователь авторизован, FromContext вернет структуру и true.
// Иначе nil и false.
func FromContext(c *gin.Context) (*User, bool) {
	st := session.FromContext(c)
	if st.GetInt("user_id") == 0 {
		return nil, false
	}
	userid := uint(st.GetInt("user_id"))
	userSes, ok := imStorage[userid]
	if !ok {
		return nil, false
	}
	return userSes.format(), true
}

// AutoLogin запишет факт авторизации в сессию пользователя.
// Он перезапишет старые данные об авторизации, если таковые имеются.
func (u *User) AutoLogin(c *gin.Context) {
	st := session.FromContext(c)
	st.SetInt("user_id", int(u.ID))
}

// Destroy стирает данные о пользователе из базы данных.
func (u *User) Destroy() {
	delete(imStorage, u.ID)
	for _, v := range u.Lists {
		v.Destroy()
	}
}
