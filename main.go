package main

import (
	"github.com/gin-gonic/gin"
	"school/todo"
//	_ "github.com/lib/pq"
)

func main() {
	r := gin.Default()
	s := todo.Todohandler {}
	r.GET("/api/todos", s.GetTodosHandler)
	r.GET("/api/todos/:id", s.GetTodosByIdHandler)
	r.POST("/api/todos", s.PostTodosHandler)
	r.DELETE("/api/todos/:id", s.DeleteTodosByIDHanderler)
	r.PUT("/api/todos/:id", s.PutTodosByIDHanderler)
	r.Run(":1234")
}

