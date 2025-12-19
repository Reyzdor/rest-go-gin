package handlers

import (
	"Application/models"
	"Application/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Login    string `form:"login"`
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

	var user *models.User
	var err error

	if strings.Contains(input.Login, "@") {
		user, err = repository.GetUserByEmail(input.Login)
	} else {
		user, err = repository.GetUserByUsername(input.Login)
	}

	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Wrong password"})
		return
	}

	c.Redirect(http.StatusFound, "/index")

}
