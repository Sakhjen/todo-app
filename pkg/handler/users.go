package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getMe(c *gin.Context) {
	id, _ := c.Get(userCtx)
	c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"id": id,
		},
	)
}
