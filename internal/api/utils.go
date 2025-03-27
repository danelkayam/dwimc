package api

import (
	api_model "dwimc/internal/api/model"
	"dwimc/internal/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func bindJsonOrErrorResponse(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		return handleErrorResponse(c, model.ErrInvalidArgs)
	}

	return false
}

func handleErrorResponse(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, model.ErrDatabase),
		errors.Is(err, model.ErrOperationFailed),
		errors.Is(err, model.ErrInternal):
		c.JSON(http.StatusInternalServerError, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Something went wrong",
			},
		})
		return true

	case errors.Is(err, model.ErrItemNotFound):
		c.JSON(http.StatusNotFound, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Not found",
			},
		})
		return true

	case errors.Is(err, model.ErrItemConflict):
		c.JSON(http.StatusConflict, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Conflict",
			},
		})
		return true

	case errors.Is(err, model.ErrInvalidArgs):
		c.JSON(http.StatusBadRequest, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Bad request",
			},
		})
		return true

	case errors.Is(err, model.ErrUnauthenticated):
		c.JSON(http.StatusUnauthorized, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Unauthenticated",
			},
		})
		return true

	case errors.Is(err, model.ErrUnauthorized):
		c.JSON(http.StatusForbidden, api_model.Response{
			Error: &api_model.ErrorResponse{
				Message: "Unauthorized",
			},
		})
		return true
	}

	return false
}
