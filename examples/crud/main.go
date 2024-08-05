package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Pla9er9/template_engine"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID      int
	Name    string
	Email   string
	IsAdult bool
}

var usersMap = make(map[int]User)
var engine = templateEngine.GetTemplateEngine()

func main() {
	r := gin.Default()

	// go SetUsers()

	r.GET("/", func(c *gin.Context) {
		users := make([]User, 0)
		for _, value := range usersMap {
			users = append(users, value)
		}

		variables := map[string]any{
			"users": users,
		}
		html, err := engine.RenderTemplateFromFile("templates/index.got", variables)
		if err != nil {
			panic(err)
		}

		c.Data(200, "text", []byte(html))
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, ok := usersMap[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		variables := map[string]any{
			"user": user,
		}
		html, err := engine.RenderTemplateFromFile("templates/user.got", variables)
		if err != nil {
			panic(err)
		}

		c.Data(200, "text", []byte(html))
	})

	r.GET("/users/:id/edit", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, ok := usersMap[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		variables := map[string]any{
			"user": user,
		}
		html, err := engine.RenderTemplateFromFile("templates/edit.got", variables)
		if err != nil {
			panic(err)
		}

		c.Data(200, "text", []byte(html))
	})

	r.GET("/new", func(c *gin.Context) {
		html, err := engine.RenderTemplateFromFile("templates/new.got", nil)
		if err != nil {
			panic(err)
		}

		c.Data(200, "text", []byte(html))
	})

	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newID := len(usersMap) + 1
		user.ID = newID
		usersMap[newID] = user

		c.JSON(http.StatusCreated, gin.H{"user": user})
	})

	r.PUT("/users/:id", func(c *gin.Context) {
		var user User
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		user.ID = id
		usersMap[id] = user

		c.JSON(http.StatusCreated, gin.H{"user": user})
	})

	r.DELETE("/users/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		delete(usersMap, id)
	})

	r.Run()
}

func SetUsers() {
	for i := 1; i <= 25; i++ {
		usersMap[i] = User{
			ID:    i,
			Name:  fmt.Sprintf("Alex-%v", i),
			Email: fmt.Sprintf("email_%v@mail.com", i),
			IsAdult: i % 2 == 0,
		}

	}
}
