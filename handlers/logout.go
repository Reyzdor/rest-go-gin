package handlers

import (
	"Application/database"
	"Application/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {
		repository.DeleteSession(database.DB, refreshToken)
	}

	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.Redirect(http.StatusFound, "/main")
}
