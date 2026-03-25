package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//- GetCartInfo godoc
//- @Summary      Get cart info
//- @Description  Returns draft ID and item count for current user
//- @Tags         Cart
//- @Security     BearerAuth
//- @Produce      json
//- @Success      200 {object} map[string]interface{} "draft_id, count"
//- @Failure      401 {object} map[string]interface{} "error"
//- @Router       /cart [get]
func (h *Handler) GetCartInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	draftID, count, err := h.Repo.GetCartInfo(userID.(uint))
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"draft_id": draftID,
		"count":    count,
	})
}

// AddItem godoc
// @Summary      Add electrolyte to cart
// @Description  Adds a electrolyte to the current user's draft calculation
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body object true "Item details" example({"electrolyte_id":1,"volume":50,"comment":"main"})
// @Success      201 {object} map[string]interface{} "message: item added"
// @Failure      400 {object} map[string]interface{} "error"
// @Failure      401 {object} map[string]interface{} "error"
// @Router       /cart/items [post]
func (h *Handler) AddItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var input struct {
		ElectrolyteID uint   `json:"electrolyte_id" binding:"required"`
		Volume        int    `json:"volume" binding:"required,min=1"`
		Comment       string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.AddItemToDraft(userID.(uint), input.ElectrolyteID, input.Volume, input.Comment); err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "item added"})
}

// UpdateItem godoc
// @Summary      Update cart item
// @Description  Updates volume and comment of a cart item
// @Tags         Cart
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Item ID"
// @Param        request body object true "Item details" example({"volume":60,"comment":"updated"})
// @Success      200 {object} map[string]interface{} "message: item updated"
// @Failure      400 {object} map[string]interface{} "error"
// @Failure      401 {object} map[string]interface{} "error"
// @Failure      404 {object} map[string]interface{} "error"
// @Router       /cart/items/{id} [put]
func (h *Handler) UpdateItem(c *gin.Context) {
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
	var input struct {
		Volume  int    `json:"volume" binding:"required,min=1"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Repo.UpdateItem(userID.(uint), uint(id), input.Volume, input.Comment); err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item updated"})
}

// DeleteItem godoc
// @Summary      Delete cart item
// @Description  Removes item from cart
// @Tags         Cart
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Item ID"
// @Success      200 {object} map[string]interface{} "message: item deleted"
// @Failure      401 {object} map[string]interface{} "error"
// @Failure      404 {object} map[string]interface{} "error"
// @Router       /cart/items/{id} [delete]
func (h *Handler) DeleteItem(c *gin.Context) {
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
	if err := h.Repo.DeleteItem(userID.(uint), uint(id)); err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}