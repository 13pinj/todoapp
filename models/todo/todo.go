package todo

import (
	"errors"

	"github.com/13pinj/todoapp/models"
	"github.com/jinzhu/gorm"
)

// Todo - это cтруктура одного "дела".
type Todo struct {
	gorm.Model
	// Статус выполнения дела.
	Done bool
	// Краткое описание дела.
	Label string
	// ID списка дел, которому принадлежит Todo
	TodoListID uint
}

// New создает новый экземпляр Todo, но не сохраняет его.
// Сохранение должно производиться с помощью функции Todo.Save().
func New(l string) *Todo {
	record := Todo{}
	record.Done = false
	record.Label = l
	return &record
}

// Find возвращает Todo, сохраненный в базе и имеющий заданный id.
// В случае если Todo не был найден, Find вернет вторым значением false.
// В случае успеха, второе возвращаемое значение будет true.
func Find(id uint) (*Todo, bool) {
	t := &Todo{}
	err := models.DB.Find(t, id).Error
	if err != nil {
		return nil, false
	}
	return t, true
}

// FindByList возвращает все Todo в базе, связанные со списком,
// имеющим заданный ID.
func FindByList(listid uint) []*Todo {
	slice := []*Todo{}
	models.DB.Where("todo_list_id = ?", listid).Find(&slice)
	return slice
}

// Save сохраняет структуру в базу. Если структура не была прежде сохранена,
// она будет добавлена в базу, в противном случае, ее представление в базе
// будет обновлено.
// Save возвращает ненулевую ошибку в случае невалидности одного их полей структуры.
// В этом случае сохранение не будет произведено.
func (t *Todo) Save() error {
	if t.Label == "" {
		return errors.New("Текст задания не должен быть пустым")
	}
	err := models.DB.Save(t).Error
	if err != nil {
		return err
	}
	return nil
}

// Destroy удаляет структуры из базы.
func (t *Todo) Destroy() {
	if models.DB.NewRecord(t) {
		return
	}
	models.DB.Delete(t)
}

func init() {
	models.DB.AutoMigrate(&Todo{})
}
