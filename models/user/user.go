package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/13pinj/todoapp/models"
	"github.com/13pinj/todoapp/models/session"
	"github.com/13pinj/todoapp/models/todolist"
)

// User - структура модели пользователя
type User struct {
	gorm.Model
	Name    string
	PwdHash string
	Lists   []*todolist.TodoList
}

// LoadLists загружает из базы списки дел пользователя в поле Lists
func (u *User) LoadLists() {
	if u.Lists != nil {
		return
	}
	u.Lists = todolist.FindByUser(u.ID)
}

func validateName(name string) bool {
	err := models.DB.Where("name = ?", name).First(&User{}).Error
	return err != nil
}

// Register добавляет нового пользователя в базу, и возвращает его структуру,
// если введенные поля корректны. В противном случае Register возвращает ошибку.
func Register(name string, password string) (u *User, errs []error) {
	ok := strings.Contains(name, " ")
	if ok {
		errs = append(errs, errors.New("Пробел в имени запрещен"))
	}
	if len([]rune(name)) < 4 {
		errs = append(errs, errors.New("Имя слишком короткое (минимум 4 символа)"))
	}
	if len([]rune(password)) < 6 {
		errs = append(errs, errors.New("Пароль слишком короткий (минимум 6 символов)"))
	}
	if !validateName(name) {
		errs = append(errs, errors.New("Имя кем-то занято"))
	}
	if errs != nil {
		return
	}
	hash := sha1.Sum([]byte(password))
	u = &User{
		Name:    name,
		PwdHash: fmt.Sprintf("%x", hash),
	}
	err := models.DB.Save(u).Error
	if err != nil {
		errs = append(errs, err)
	}
	return
}

// Login выполняет авторизацию пользователей.
// Если введенные имя и пароль действительны, Login запишет факт авторизации
// в сессию пользователя и вернет первым значением структуру пользователя,
// а вторым true. В противном случае - nil и false.
// Login перезапишет старые данные об авторизации, если таковые имеются.
func Login(c *gin.Context, name string, password string) (*User, bool) {
	user := &User{}
	err := models.DB.Where("name = ?", name).First(user).Error
	if err != nil {
		return nil, false
	}
	hash := sha1.Sum([]byte(password))
	str := fmt.Sprintf("%x", hash)
	if str != user.PwdHash {
		return nil, false
	}
	st := session.FromContext(c)
	st.SetInt("user_id", int(user.ID))
	return user, true
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
	user := &User{}
	err := models.DB.Find(user, userid).Error
	if err != nil {
		return nil, false
	}
	return user, true
}

// AutoLogin запишет факт авторизации в сессию пользователя.
// Он перезапишет старые данные об авторизации, если таковые имеются.
func (u *User) AutoLogin(c *gin.Context) {
	st := session.FromContext(c)
	st.SetInt("user_id", int(u.ID))
}

// Destroy стирает данные о пользователе из базы данных.
func (u *User) Destroy() {
	models.DB.Delete(u)
}

func init() {
	models.DB.AutoMigrate(&User{})
}
