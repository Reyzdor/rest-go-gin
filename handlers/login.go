package handlers

import (
	"Application/auth"
	"Application/database"
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

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Email, user.Username)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Server error (A)"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Server error (R)"})
		return
	}

	err = repository.SaveSession(database.DB, user.ID, refreshToken)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{"error": "Server error (E)"})
		return
	}

	c.SetCookie("auth_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.Redirect(http.StatusFound, "/main")

}
