package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetConcentrations(c *gin.Context) {
	status := c.Query("status")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	var fromDate, toDate *time.Time
	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err == nil {
			fromDate = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err == nil {
			toDate = &t
		}
	}
	concs, err := h.Repo.GetUserConcentrations(status, fromDate, toDate)
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": concs})
}

func (h *Handler) GetConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	conc, items, err := h.Repo.GetConcentrationByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	if conc.Status == "удалён" {
		c.JSON(http.StatusNotFound, gin.H{"error": "concentration deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"concentration": conc,
		"items":         items,
	})
}

func (h *Handler) UpdateConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input struct {
		Result string `json:"result"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conc, _, err := h.Repo.GetConcentrationByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	if conc.Status != "черновик" && conc.Status != "сформирован" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update in this status"})
		return
	}
	conc.Result = input.Result
	if err := h.Repo.UpdateConcentration(conc); err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conc})
}

func (h *Handler) FormConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Repo.SetStatusFormed(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration formed"})
}

func (h *Handler) FinishConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Repo.SetStatusFinished(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration finished"})
}

func (h *Handler) RejectConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Repo.SetStatusRejected(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration rejected"})
}

func (h *Handler) DeleteConcentration(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.Repo.DeleteConcentration(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration deleted"})
}