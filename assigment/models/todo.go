package models

type Todo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UpdateTodo struct {
	Name string `json:"name"`
}
