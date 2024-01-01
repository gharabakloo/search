package http

import (
	"errors"

	"github.com/gin-gonic/gin"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/internal/service"
	"gharabakloo/search/pkg/myerr"
)

type Handler struct {
	cfg     *entity.Config
	service service.SearchService
}

func NewHandler(
	cfg *entity.Config,
	service service.SearchService,
) *Handler {
	return &Handler{
		cfg:     cfg,
		service: service,
	}
}

func (h *Handler) sendError(c *gin.Context, statusCode int, err error) {
	keyError := "error"
	var e *myerr.Error
	if errors.As(err, &e) {
		c.JSON(statusCode, gin.H{
			keyError: *e,
		})
		return
	}

	c.JSON(statusCode, gin.H{
		keyError: err.Error(),
	})
}
