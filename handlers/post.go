package handlers

import (
	"Application/database"
	"Application/models"
	"Application/repository"
	"database/sql"
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

	if _, err := os.Stat("./static/uploads"); os.IsNotExist(err) {
		os.MkdirAll("./static/uploads", 0755)
	}

	result, err := database.DB.Exec(
		"INSERT INTO posts (title, content, price, user_id) VALUES(?,?,?,?)",
		title, content, price, userID,
	)

	if err != nil {
		c.String(500, "Error creating post: %v", err)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		c.String(500, "Error getting post ID: %v", err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.String(400, "Error getting form: %v", err)
		return
	}

	files := form.File["images[]"]
	mainImageSet := false

	for i, file := range files {
		filename := fmt.Sprintf("post_%d_%d_%s", postID, i, file.Filename)
		filePath := filepath.Join("./static/uploads", filename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			fmt.Printf("Error saving file %s: %v\n", file.Filename, err)
			continue
		}

		imagePath := "/static/uploads/" + filename

		isMain := false
		if !mainImageSet {
			isMain = true
			mainImageSet = true

			database.DB.Exec(
				"UPDATE posts SET main_image = ? WHERE id = ?",
				imagePath, postID,
			)
		}

		_, err := database.DB.Exec(
			"INSERT INTO post_images (post_id, image_path, is_main, sort_order) VALUES(?,?,?,?)",
			postID, imagePath, isMain, i,
		)

		if err != nil {
			fmt.Printf("Error saving image to DB: %v\n", err)
		}
	}

	c.Redirect(302, "/admin")
}

func ToursPage(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT p.*, u.username
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY created_at DESC
	`)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Error loading tours",
		})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var mainImage sql.NullString

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&mainImage,
			&post.Price,
			&post.UserID,
			&post.CreatedAt,
			&post.Username,
		)

		if err != nil {
			fmt.Printf("Error scanning post: %v\n", err)
			continue
		}

		if mainImage.Valid {
			post.MainImage = mainImage.String
		}

		posts = append(posts, post)
	}

	c.HTML(http.StatusOK, "tours.html", gin.H{
		"posts": posts,
	})
}
