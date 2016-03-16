package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/13pinj/todoapp/models"
	"github.com/13pinj/todoapp/models/session"
	"github.com/13pinj/todoapp/models/todolist"
)

// Роли пользователя.
// Возможные значения User.Role.
const (
	DefaultRole = ""
	AdminRole   = "admin"
)

// User - структура модели пользователя
type User struct {
	gorm.Model
	Name      string
	PwdHash   string
	Role      string
	Lists     []*todolist.TodoList `gorm:"-"`
	VisitedAt time.Time
}

func hashPwd(pwd string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(pwd)))
}

// Register добавляет нового пользователя в базу, и возвращает его структуру,
// если введенные поля корректны. В противном случае Register возвращает ошибку.
func Register(name string, password string) (u *User, errs []error) {
	if strings.Contains(name, " ") {
		errs = append(errs, errors.New("Пробел в имени запрещен"))
	}
	if len([]rune(name)) < 4 {
		errs = append(errs, errors.New("Имя слишком короткое (минимум 4 символа)"))
	}
	if len([]rune(password)) < 6 {
		errs = append(errs, errors.New("Пароль слишком короткий (минимум 6 символов)"))
	}
	if _, ok := Find(name); ok {
		errs = append(errs, errors.New("Имя кем-то занято"))
	}
	if errs != nil {
		return
	}
	u = &User{
		Name:    name,
		PwdHash: hashPwd(password),
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
	if hashPwd(password) != user.PwdHash {
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
	user.MarkVisit()
	return user, true
}

// Find находит пользователя в базе по указанному имени.
// Второе возвращаемое значение будет равно false в случае безуспешного поиска.
func Find(name string) (*User, bool) {
	user := &User{}
	err := models.DB.Where("name = ?", name).First(user).Error
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

// MarkVisit обновляет поле VisitedAt и сохраняет в базу.
func (u *User) MarkVisit() {
	u.VisitedAt = time.Now()
	models.DB.Save(u)
}

// LoadLists загружает из базы списки дел пользователя в поле Lists
func (u *User) LoadLists() {
	if u.Lists != nil {
		return
	}
	u.Lists = todolist.FindByUser(u.ID)
}

// SetRole задает пользователю новую роль и сохраняют в базу.
func (u *User) SetRole(r string) {
	u.Role = r
	models.DB.Save(u)
}

// Admin возвращает true, если пользователь относится к администрации.
func (u *User) Admin() bool {
	return u.Role == AdminRole
}

// Destroy стирает данные о пользователе из базы данных.
func (u *User) Destroy() {
	models.DB.Delete(u)
}

// Count возвращает общее количество всех существующих пользователей.
func Count() (n int) {
	models.DB.Model(&User{}).Count(&n)
	return
}

// Pages возвращает количество страниц, на которые мог бы поместиться
// список всех существующих пользователей по n элементов на страницу.
func Pages(n int) int {
	return (Count()-1)/n + 1
}

// SortMode определяет режим сортировки выборки пользователей.
type SortMode string

// Возможные варианты сортировки выборок пользователей.
const (
	ByID            SortMode = "id"
	ByName          SortMode = "name"
	ByCreatedAt     SortMode = "created_at"
	ByIDDesc        SortMode = "id desc"
	ByNameDesc      SortMode = "name desc"
	ByCreatedAtDesc SortMode = "created_at desc"
)

// FindPage возвращает список пользователей на i-й странице, если бы они
// размещались по n штук на страницу, отсортированные по sortBy.
// Отсчет страниц ведется с единицы.
func FindPage(i, n int, sortBy SortMode) (us []*User) {
	models.DB.Limit(n).Model(&User{}).Offset(n * (i - 1)).Order(string(sortBy)).Find(&us)
	return
}

var initUser = &User{
	Name:    "root",
	PwdHash: hashPwd("12345678"),
	Role:    AdminRole,
}

func init() {
	initializeUsers()
}

func initializeUsers() {
	dummy := &User{}
	if !models.DB.HasTable(dummy) {
		err := models.DB.CreateTable(dummy).Error
		if err != nil {
			panic(err)
		}
		err = models.DB.Create(initUser).Error
		if err != nil {
			panic(err)
		}
	}
	models.DB.AutoMigrate(dummy)
}
