package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCartInfo(c *gin.Context) {
	draftID, count, err := h.Repo.GetCartInfo()
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"draft_id": draftID,
		"count":    count,
	})
}