package dto

type ResponseDTO struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ArticleResDTO struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}
