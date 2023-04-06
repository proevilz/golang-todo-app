package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"proevilz/api/db"
	"proevilz/api/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getTodos(c *gin.Context) {
	var todos []models.Todo
	db.DB.Find(&todos)
	c.IndentedJSON(http.StatusOK, todos)
}

func updateTodo(c *gin.Context) {
	type TodoUpdate struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	id := c.Param("id")
	var todo models.Todo
	result := db.DB.First(&todo, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	} else if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var update TodoUpdate
	err := c.BindJSON(&update)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if update.Title != nil {
		todo.Title = *update.Title
	}
	if update.Completed != nil {
		todo.Completed = *update.Completed
	}

	result = db.DB.Updates(&todo)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, todo)

}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	todo := models.Todo{ID: id}
	result := db.DB.Delete(&todo)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.IndentedJSON(http.StatusNoContent, id)
}

func createTodo(c *gin.Context) {
	var todo models.Todo
	err := json.NewDecoder(c.Request.Body).Decode(&todo)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTodo := models.Todo{ID: uuid.New().String(), Title: todo.Title, Completed: false}

	result := db.DB.Create(&newTodo)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.IndentedJSON(http.StatusCreated, newTodo)

}
func main() {

	db.ConnectDB()
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/todos", getTodos)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.POST("/todos", createTodo)
	router.Run("localhost:8080")
}
