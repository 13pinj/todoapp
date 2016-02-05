package todolist

import (
	"errors"

	"github.com/13pinj/todoapp/models/todo"
)

// TodoList - структура списка дел
type TodoList struct {
	ID uint
	// Заголовок списка
	Title string
	// Список todo
	Todos []*todo.Todo
	// ID пользователя, которому принадлежит список
	UserID uint
}

// Представление TodoList в памяти
type imTodoList struct {
	ID    uint
	Title string
}

// Хранилище списков дел внутри памяти.
var imStorage = make(map[uint]*imTodoList)
var key uint

// Функция форматирования внутреннего представления TodoList во внешее представление.
func (l *imTodoList) format() *TodoList {
	return &TodoList{
		ID:    l.ID,
		Title: l.Title,
		Todos: todo.FindByList(l.ID),
	}
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
	return nil
}

// Find возвращает TodoList, сохраненный в базе и имеющий заданный id.
// В случае если TodoList не был найден, Find вернет вторым значением false.
// В случае успеха, второе возвращаемое значение будет true.
func Find(id uint) (*TodoList, bool) {
	record, ok := imStorage[id]
	if ok {
		return record.format(), true
	}
	return nil, false
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
	if l.ID == 0 {
		key++
		l.ID = key
		imStorage[l.ID] = &imTodoList{
			ID:    l.ID,
			Title: l.Title,
		}
	} else {
		imStorage[l.ID].Title = l.Title
	}
	return nil
}

// Destroy удаляет структуры из базы.
func (l *TodoList) Destroy() {
	delete(imStorage, l.ID)
	for _, v := range l.Todos {
		v.Destroy()
	}
}

// Len возвращает количество всех дел в списке.
func (l *TodoList) Len() int {
	return len(l.Todos)
}

// LenUndone возвращает количество незавершенных дел в списке.
func (l *TodoList) LenUndone() int {
	k := 0
	for _, v := range l.Todos {
		if !v.Done {
			k++
		}
	}
	return k
}

// LenDone возвращает количество завершенных дел в списке.
func (l *TodoList) LenDone() int {
	k := 0
	for _, v := range l.Todos {
		if v.Done {
			k++
		}
	}
	return k
}

// Add добавляет в список новое назавершенное дело c текстом lbl и сохраняет его в базу.
// Add не позволяет добавлять дела в несохраненную в базу списки, о чем сообщает ошибкой.
// Add также возвращает ошибки сохранения Todo.
func (l *TodoList) Add(lbl string) error {
	if l.ID == 0 {
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
	for _, v := range l.Todos {
		if !v.Done {
			slice = append(slice, v)
		}
	}
	return slice
}

// Done возвращает список завершенных дел.
func (l *TodoList) Done() []*todo.Todo {
	slice := []*todo.Todo{}
	for _, v := range l.Todos {
		if v.Done {
			slice = append(slice, v)
		}
	}
	return slice
}
