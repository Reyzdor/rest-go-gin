package main

import (
	"Application/database"
	"Application/handlers"
	"Application/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitSQLite(); err != nil {
		log.Fatalf("database init error: %v", err)
	}
	defer database.DB.Close()

	r := gin.Default()

	r.Use(middleware.AuthMiddleware())

	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", []byte{})
	})

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})
	r.GET("/tours", handlers.ToursPage)

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})

	r.GET("/main", func(c *gin.Context) {
		isAuth, _ := c.Get("is_authenticated")
		email, _ := c.Get("user_email")
		username, _ := c.Get("username")

		data := gin.H{
			"Title": "Travelo - Главная",
		}

		if isAuthBool, ok := isAuth.(bool); ok && isAuthBool {
			data["Authenticated"] = true
			data["Email"] = email
			data["Username"] = username
		} else {
			data["Authenticated"] = false
		}

		c.HTML(200, "index.html", data)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/main")
	})

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)
	r.GET("/admin", handlers.AdminPage)
	r.POST("/create-post", handlers.CreatePost)

	r.Run()
}
