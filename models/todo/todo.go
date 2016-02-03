package todo

import "errors"

// Todo - это cтруктура одного "дела".
type Todo struct {
	ID uint
	// Статус выполнения дела.
	Done bool
	// Краткое описание дела.
	Label string
}

var st = make(map[uint]*Todo)
var key uint

// New создает новый экземпляр Todo, но не сохраняет его.
// Сохранение должно производиться с помощью функции Todo.Save().
func New(l string) *Todo {
	record := Todo{}
	record.Done = false
	record.Label = l
	return &record
}

func stCopy(t *Todo) *Todo {
	return &Todo{
		ID:    t.ID,
		Done:  t.Done,
		Label: t.Label,
	}
}

// Find возвращает Todo, сохраненный в базе и имеющий заданный id.
// В случае если Todo не был найден, Find вернет вторым значением false.
// В случае успеха, второе возвращаемое значение будет true.
func Find(id uint) (*Todo, bool) {
	record, ok := st[id]
	if ok {
		return stCopy(record), true
	}
	return nil, false
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
	if t.ID == 0 {
		key++
		t.ID = key
		st[t.ID] = stCopy(t)
	} else {
		st[t.ID].Done = t.Done
		st[t.ID].Label = t.Label
	}
	return nil
}

// Destroy удаляет структуры из базы.
func (t *Todo) Destroy() {
	delete(st, t.ID)
}
