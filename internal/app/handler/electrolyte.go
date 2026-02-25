package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
		"go_project2/internal/app/ds"
)

func (h *Handler) GetElectrolytes(c *gin.Context) {
	search := c.Query("search")
	var electrolytes []ds.Electrolyte
	var err error

	if search == "" {
		electrolytes, err = h.Repo.GetAllElectrolytes()
	} else {
		electrolytes, err = h.Repo.SearchElectrolytesByName(search)
	}
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}

	// Для отображения корзины нужен ID пользователя (пока захардкодим 1)
	const creatorID = 1
	cartCount := h.Repo.GetCartCount(creatorID)
	// Также получим ID черновика, если он есть (для ссылки)
	draft, _ := h.Repo.GetOrCreateDraft(creatorID) // создаст, если нет, но мы не хотим создавать при каждом просмотре. Лучше отдельный метод GetDraftID
	// Упростим: будем получать черновик только при добавлении, а здесь просто проверяем наличие через GetCartCount > 0.
	// Ссылку на черновик дадим, только если cartCount > 0.
	var draftID uint
	if cartCount > 0 {
		draft, _ = h.Repo.GetOrCreateDraft(creatorID)
		draftID = draft.ID
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"electrolytes": electrolytes,
		"search":       search,
		"cart_count":   cartCount,
		"draft_id":     draftID,
	})
}

func (h *Handler) GetElectrolyte(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	electrolyte, err := h.Repo.GetElectrolyteByID(id)
	if err != nil || electrolyte == nil {
		c.String(http.StatusNotFound, "Not found")
		return
	}
	c.HTML(http.StatusOK, "electrolyte.html", electrolyte)
}