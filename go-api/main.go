package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: uuid.New().String(), Title: "Write presentation", Completed: false},
	{ID: uuid.New().String(), Title: "Make a cup of tea bruv", Completed: true},
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func updateTodo(c *gin.Context) {
	id := c.Param("id")

	updates := make(map[string]interface{})
	err := json.NewDecoder(c.Request.Body).Decode(&updates)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, t := range todos {
		if t.ID == id {
			updatedTodo := t
			for key, value := range updates {
				switch strings.ToLower(key) {
				case "title":
					updatedTodo.Title = value.(string)
				case "completed":
					updatedTodo.Completed = value.(bool)
				default:
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid field(s)"})
					return
				}
			}
			todos[i] = updatedTodo
			c.IndentedJSON(http.StatusOK, updatedTodo)
			return
		}
	}
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted"})
			return
		}
	}
}

func createTodo(c *gin.Context) {
	var todo Todo
	err := json.NewDecoder(c.Request.Body).Decode(&todo)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todo.ID = uuid.New().String()
	todos = append(todos, todo)
	c.IndentedJSON(http.StatusCreated, todo)

}
func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/todos", getTodos)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.POST("/todos", createTodo)
	router.Run("localhost:8080")
}
