package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/internal/handler/mapper"
)

func (h *Handler) Search(c *gin.Context) {
	key := c.Query(entity.KeySearch)
	page := entity.Page{}
	page.Number = c.Query(entity.KeyPageNumber)
	page.Size = c.Query(entity.KeyPageSize)

	books, err := h.service.Search(c.Request.Context(), key, page)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ToBooksResp(*books))
}
