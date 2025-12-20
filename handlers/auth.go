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

type RegisterInput struct {
	Username        string `form:"username" json:"username"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid input"})
		return
	}

	if !isValidEmail(input.Email) {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid email"})
		return
	}

	if !isValidPassword(input.Password) {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	if !isValidUsername(input.Username) {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Incorrect username"})
		return
	}

	if input.Password != input.ConfirmPassword {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Password do not match"})
		return
	}

	emailExists, err := repository.CheckEmailExists(input.Email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Internal serval error email"})
		return
	}

	if emailExists {
		c.HTML(http.StatusConflict, "register.html", gin.H{"error": "Email already registered"})
		return
	}

	usernameExists, err := repository.CheckUsernameExists(input.Username)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Internal server error username"})
		return
	}

	if usernameExists {
		c.HTML(http.StatusConflict, "register.html", gin.H{"error": "Username already registered"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{"error": "Hash error"})
		return
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := repository.CreateUser(user); err != nil {
		c.HTML(http.StatusConflict, "register.html", gin.H{"error": "User already exists"})
		return
	}

	createdUser, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	accessToken, err := auth.GenerateAccessToken(createdUser.ID, createdUser.Email, createdUser.Username)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	repository.SaveSession(database.DB, createdUser.ID, refreshToken)

	c.SetCookie("auth_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24*7, "/", "", false, true)

	c.Redirect(http.StatusFound, "/main")
}

func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 || at == len(email)-1 {
		return false
	}

	dot := strings.LastIndex(email, ".")
	if dot < at+2 || dot == len(email)-1 {
		return false
	}

	return true
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	return true

}

func isValidUsername(username string) bool {
	if len(username) < 4 {
		return false
	}

	for _, char := range username {
		if char < '0' || char > '9' {
			return true
		}
	}

	return false
}
