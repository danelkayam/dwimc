package api_utils

import (
	api_model "dwimc/internal/api/model"
	"dwimc/internal/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func BindJsonOrErrorResponse(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		return HandleErrorResponse(c, model.ErrInvalidArgs)
	}

	return false
}

func HandleErrorResponse(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, model.ErrDatabase),
		errors.Is(err, model.ErrOperationFailed),
		errors.Is(err, model.ErrInternal):
		log.Error().
			Err(err).
			Msg("something went wrong")

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Something went wrong",
				},
			},
		)
		return true

	case errors.Is(err, model.ErrItemNotFound):
		c.AbortWithStatusJSON(
			http.StatusNotFound,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Not found",
				},
			},
		)
		return true

	case errors.Is(err, model.ErrItemConflict):
		c.AbortWithStatusJSON(
			http.StatusConflict,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Conflict",
				},
			},
		)
		return true

	case errors.Is(err, model.ErrInvalidArgs):
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Bad request",
				},
			},
		)
		return true

	case errors.Is(err, model.ErrUnauthenticated):
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Unauthenticated",
				},
			},
		)
		return true

	case errors.Is(err, model.ErrUnauthorized):
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			api_model.Response{
				Error: &api_model.ErrorResponse{
					Message: "Unauthorized",
				},
			},
		)
		return true
	}

	return false
}
