package handlers

import (
	"Application/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Username string `form:"username_login"`
	Email    string `form:"email_login"`
	Password string `form:"password_login"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"error": "Invalid input, please fill all fields",
		})
		return
	}

	user, err := repository.GetUserByUsername(input.Username)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Wrong password"})
		return
	}

	c.HTML(http.StatusOK, "/", gin.H{"username": user.Username})

}
