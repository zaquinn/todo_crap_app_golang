package todo

type CreateTodoDTO struct {
	Id      *int   `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}
