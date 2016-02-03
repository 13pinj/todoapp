package todo_list

import "github.com/13pinj/todoapp/models/todo"

// TodoList - структура списка дел
type TodoList struct {
	ID uint
	// Заголовок списка
	Title string
	// Список todo
	Todos []*todo.Todo
}

// Представление TodoList в памяти
type imTodoList struct {
	ID    uint
	Title string
}

// Хранилище списков дел внутри памяти.
var imStorage = make(map[uint]*imTodoList)

// Функция форматирования внутреннего представления TodoList во внешее представление.
func (l *imTodoList) format() *TodoList {
	return nil
}

// New создает новый экземпляр TodoList, но не сохраняет его.
// Сохранение должно производиться с помощью функции TodoList.Save().
func New(t string) *TodoList {
	return nil
}

// Find возвращает TodoList, сохраненный в базе и имеющий заданный id.
// В случае если TodoList не был найден, Find вернет вторым значением false.
// В случае успеха, второе возвращаемое значение будет true.
func Find(id uint) (*TodoList, bool) {
	return nil, false
}

// Save сохраняет структуру в базу. Если структура не была прежде сохранена,
// она будет добавлена в базу, в противном случае, ее представление в базе
// будет обновлено.
// Save возвращает ненулевую ошибку в случае невалидности одного их полей структуры.
// В этом случае сохранение не будет произведено.
func (l *TodoList) Save() error {
	return nil
}

// Destroy удаляет структуры из базы.
func (l *TodoList) Destroy() {

}

// Len возвращает количество всех дел в списке.
func (l *TodoList) Len() int {
	return 0
}

// LenUndone возвращает количество незавершенных дел в списке.
func (l *TodoList) LenUndone() int {
	return 0
}

// LenDone возвращает количество завершенных дел в списке.
func (l *TodoList) LenDone() int {
	return 0
}

// Add добавляет в список новое назавершенное дело c текстом lbl и сохраняет его в базу.
// Add не позволяет добавлять дела в несохраненную в базу списки, о чем сообщает ошибкой.
// Add также возвращает ошибки сохранения Todo.
func (l *TodoList) Add(lbl string) error {
	return nil
}

// Undone возвращает список незавершенных дел.
func (l *TodoList) Undone() []*todo.Todo {
	return nil
}

// Done возвращает список завершенных дел.
func (l *TodoList) Done() []*todo.Todo {
	return nil
}
