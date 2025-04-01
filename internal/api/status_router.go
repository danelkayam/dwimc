package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusRouter struct {

}

func NewStatusRouter() *StatusRouter {
	return &StatusRouter{}
}

func (r *StatusRouter) Health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string {
		"status": "ok",
	})
}

func (r *StatusRouter) Live(c *gin.Context) {
	c.Status(http.StatusOK)
}