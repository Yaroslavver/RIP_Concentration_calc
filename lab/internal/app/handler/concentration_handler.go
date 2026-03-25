package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetConcentrations godoc
// @Summary      Get user's concentrations
// @Description  Returns list of user's concentrations (draft and deleted excluded)
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        status query string false "Filter by status"
// @Param        from query string false "Filter by date from (YYYY-MM-DD)"
// @Param        to query string false "Filter by date to (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Router       /concentrations [get]
func (h *Handler) GetConcentrations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	concs, err := h.Repo.GetUserConcentrations(userID.(uint), status, fromDate, toDate)
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": concs})
}

// GetAllConcentrationsForModerator godoc
// @Summary      Get all concentrations (moderator only)
// @Description  Returns all concentrations from all users (draft and deleted excluded)
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        status query string false "Filter by status"
// @Param        from query string false "Filter by date from (YYYY-MM-DD)"
// @Param        to query string false "Filter by date to (YYYY-MM-DD)"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Router       /concentrations/all [get]
func (h *Handler) GetAllConcentrationsForModerator(c *gin.Context) {
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
	concs, err := h.Repo.GetAllConcentrations(status, fromDate, toDate)
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": concs})
}

// GetConcentration godoc
// @Summary      Get concentration by ID
// @Description  Returns detailed information about a specific concentration with its items
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id} [get]
func (h *Handler) GetConcentration(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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
	// Проверяем, что пользователь имеет доступ к этой заявке (создатель или модератор)
	isMod, _ := c.Get("isModerator")
	if conc.CreatorID != userID.(uint) && !isMod.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"concentration": conc,
		"items":         items,
	})
}

// UpdateConcentration godoc
// @Summary      Update concentration
// @Description  Updates result and description fields of a concentration
// @Tags         Concentrations
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Param        request body object true "Update data" example({"result":"[H⁺] = 0.045 моль/л","description":"Mixture description"})
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id} [put]
func (h *Handler) UpdateConcentration(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	conc, _, err := h.Repo.GetConcentrationByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	// Проверяем права
	if conc.CreatorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only creator can update"})
		return
	}
	if conc.Status != "черновик" && conc.Status != "сформирован" {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update in this status"})
		return
	}
	var input struct {
		Result      string `json:"result"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conc.Result = input.Result
	conc.Description = input.Description
	if err := h.Repo.UpdateConcentration(conc); err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": conc})
}

// FormConcentration godoc
// @Summary      Form concentration (user action)
// @Description  Changes status from draft to formed and calculates result
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id}/formed [put]
func (h *Handler) FormConcentration(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	conc, _, err := h.Repo.GetConcentrationByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	if conc.CreatorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only creator can form"})
		return
	}
	if err := h.Repo.SetStatusFormed(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration formed"})
}

// FinishConcentration godoc
// @Summary      Finish concentration (moderator action)
// @Description  Changes status to finished, sets moderator and finish date
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id}/finish [put]
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

// RejectConcentration godoc
// @Summary      Reject concentration (moderator action)
// @Description  Changes status to rejected, sets moderator and finish date
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id}/reject [put]
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

// DeleteConcentration godoc
// @Summary      Delete concentration (soft delete)
// @Description  Changes status to deleted (logical deletion)
// @Tags         Concentrations
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Concentration ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]interface{}
// @Failure      401 {object} map[string]interface{}
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Router       /concentrations/{id} [delete]
func (h *Handler) DeleteConcentration(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	conc, _, err := h.Repo.GetConcentrationByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	if conc.CreatorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only creator can delete"})
		return
	}
	if err := h.Repo.DeleteConcentration(uint(id)); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "concentration deleted"})
}