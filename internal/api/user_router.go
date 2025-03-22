package api

import (
	"dwimc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /api/user/auth/register - creates new user
// POST /api/user/auth/signin - get new user_token (JWT) for existing user

// Requires: user_token
// GET     /api/user/me - get current user based on user_token
// PUT     /api/user/ - update current user profile (name)
// DELETE  /api/user/ - delete current user, will delete all devices & locations

type UserRouter struct {
	userService services.UserService
}

func NewUserRouter(userService services.UserService) *UserRouter {
	return &UserRouter{userService: userService}
}

func (r *UserRouter) RegisterUser(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *UserRouter) SignInUser(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *UserRouter) GetCurrentUser(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *UserRouter) UpdateCurrentUser(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *UserRouter) DeleteCurrentUser(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}
