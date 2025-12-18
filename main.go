package main

import (
	"Application/database"
	"Application/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitSQLite(); err != nil {
		log.Fatalf("database init error: %v", err)
	}
	defer database.DB.Close()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/main", func(c *gin.Context) {
		c.HTML(200, "main.html", nil)
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.Run()

}
