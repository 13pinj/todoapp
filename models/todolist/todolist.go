package todolist

import (
	"errors"
	"fmt"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/13pinj/todoapp/models"
	"github.com/13pinj/todoapp/models/todo"
)

// TodoList - структура списка дел
type TodoList struct {
	gorm.Model
	// Заголовок списка
	Title string
	// Список todo
	Todos []*todo.Todo
	// ID пользователя, которому принадлежит список
	UserID uint
}

// LoadTodos загружает из базы задания списка в поле Todos
func (l *TodoList) LoadTodos() {
	if l.Todos != nil {
		return
	}
	l.Todos = todo.FindByList(l.ID)
}

// New создает новый экземпляр TodoList, но не сохраняет его.
// Сохранение должно производиться с помощью функции TodoList.Save().
func New(t string) *TodoList {
	record := TodoList{}
	record.Title = t
	return &record
}

// FindByUser возвращает все списки в базе, принадлежащие пользователю
// с заданным ID.
func FindByUser(userID uint) []*TodoList {
	slice := []*TodoList{}
	models.DB.Where("user_id = ?", userID).Order("created_at ASC").Find(&slice)
	return slice
}

// Find возвращает TodoList, сохраненный в базе и имеющий заданный id.
// В случае если TodoList не был найден, Find вернет вторым значением false.
// В случае успеха, второе возвращаемое значение будет true.
func Find(id uint) (*TodoList, bool) {
	tl := &TodoList{}
	err := models.DB.Find(tl, id).Error
	if err != nil {
		return nil, false
	}
	return tl, true
}

// Save сохраняет структуру в базу. Если структура не была прежде сохранена,
// она будет добавлена в базу, в противном случае, ее представление в базе
// будет обновлено.
// Save возвращает ненулевую ошибку в случае невалидности одного их полей структуры.
// В этом случае сохранение не будет произведено.
func (l *TodoList) Save() error {
	if l.Title == "" {
		return errors.New("Заголовок списка не должен быть пустым")
	}
	err := models.DB.Save(l).Error
	if err != nil {
		return err
	}
	return nil
}

// Destroy удаляет структуры из базы.
func (l *TodoList) Destroy() {
	for _, v := range l.Todos {
		v.Destroy()
	}
	if models.DB.NewRecord(l) {
		return
	}
	models.DB.Delete(l)
}

// Len возвращает количество всех дел в списке.
func (l *TodoList) Len() int {
	var count int
	models.DB.Model(&todo.Todo{}).Where("todo_list_id = ?", l.ID).Count(&count)
	return count
}

// LenUndone возвращает количество незавершенных дел в списке.
func (l *TodoList) LenUndone() int {
	var count int
	models.DB.Model(&todo.Todo{}).Where("todo_list_id = ? AND done = ?", l.ID, false).Count(&count)
	return count
}

// LenDone возвращает количество завершенных дел в списке.
func (l *TodoList) LenDone() int {
	var count int
	models.DB.Model(&todo.Todo{}).Where("todo_list_id = ? AND done = ?", l.ID, true).Count(&count)
	return count
}

// Add добавляет в список новое назавершенное дело c текстом lbl и сохраняет его в базу.
// Add не позволяет добавлять дела в несохраненную в базу списки, о чем сообщает ошибкой.
// Add также возвращает ошибки сохранения Todo.
func (l *TodoList) Add(lbl string) error {
	if models.DB.NewRecord(l) {
		return errors.New("Нельзя добавлять дела в несохраненный список")
	}
	sd := todo.New(lbl)
	sd.TodoListID = l.ID
	err := sd.Save()
	if err != nil {
		return err
	}
	l.Todos = append(l.Todos, sd)
	return nil
}

// Undone возвращает список незавершенных дел.
func (l *TodoList) Undone() []*todo.Todo {
	slice := []*todo.Todo{}
	models.DB.Model(&todo.Todo{}).Where("todo_list_id = ? AND done = ?", l.ID, false).Find(&slice)
	return slice
}

// Done возвращает список завершенных дел.
func (l *TodoList) Done() []*todo.Todo {
	slice := []*todo.Todo{}
	models.DB.Model(&todo.Todo{}).Where("todo_list_id = ? AND done = ?", l.ID, true).Find(&slice)
	return slice
}

// Path возвращает путь к странице списка.
func (l *TodoList) Path() string {
	return fmt.Sprintf("/list/%d", l.ID)
}

func init() {
	models.DB.AutoMigrate(&TodoList{})
}
