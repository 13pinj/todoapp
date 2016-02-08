package todos

import "github.com/gin-gonic/gin"

// Index выводит страницу со всеми списками дел текущего пользователя.
// Если пользователь незалогинен, перенаправляет на главную.
// GET /
func Index(c *gin.Context) {

}

// CreateList создает новый список дел с заголовком из POST-параметра title
// и перенаправляет на страницу этого списка.
// POST /list-create
func CreateList(c *gin.Context) {

}

// ShowList выводит страницу списка дел, на которой отображается его заголовок
// и содержание.
// GET /list/:id
func ShowList(c *gin.Context) {

}

// UpdateList изменяет заголовок списка на тот, который был получен
// POST-параметре title и перенаправляет на страницу этого списка.
// POST /list/:id/update
func UpdateList(c *gin.Context) {

}

// DestroyList стирает список из базы и перенаправляет на главную.
// POST /list/:id/destroy
func DestroyList(c *gin.Context) {

}

// CreateTask создает новое задание в списке с текстом
// из POST-параметра label и перенаправляет на страницу списка.
// POST /list/:id/task-create
func CreateTask(c *gin.Context) {

}

// UpdateTask изменяет поля задания используя POST-параметры done и label.
// Поля, для который не заданы значения в параметрах запроса, должны остаться
// неизменными.
// После выполнения запроса UpdateTask перенаправляет клиент на страницу списка.
// POST /list/:id/task-update/:task-id
func UpdateTask(c *gin.Context) {

}

// DestroyTask стирает задание из списка и перенаправляет на страницу списка.
// POST /list/:id/task-destroy/:task-id
func DestroyTask(c *gin.Context) {

}
