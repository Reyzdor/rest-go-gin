package handlers

import (
	"Application/database"
	"Application/repository"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminPage(c *gin.Context) {
	userID := 1

	user, err := repository.GetUserByID(userID)
	if err != nil || user.Role != "admin" {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "You're not admin!",
		})
		return
	}

	c.HTML(http.StatusOK, "admin.html", nil)
}

func CreatePost(c *gin.Context) {
	userID := 1

	user, err := repository.GetUserByID(userID)
	if err != nil || user.Role != "admin" {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"error": "Only for admins!",
		})

		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	priceStr := c.PostForm("price")

	var price int

	if priceStr != "" {
		price, _ = strconv.Atoi(priceStr)
	}

	var imagePath string
	file, err := c.FormFile("image")
	if err == nil {
		os.MkdirAll("./static/uploads", 0755)

		filename := fmt.Sprintf("%d_%s", userID, file.Filename)
		filepath := filepath.Join("./static/uploads", filename)

		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.String(http.StatusInternalServerError, "Error saving image")
			return
		}

		imagePath = "/static/uploads/" + filename
	}

	_, err = database.DB.Exec(
		"INSERT INTO posts (title, content, image, price, user_id) VALUES(?,?,?,?,?)",
		title, content, imagePath, price, userID,
	)

	if err != nil {
		c.String(500, "Error: %v", err)
		return
	}

	c.Redirect(302, "/admin")
}
