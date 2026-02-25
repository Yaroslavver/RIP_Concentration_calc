package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCalculation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
		return
	}
	calc, items, err := h.Repo.GetCalculationByID(id)
	if err != nil {
		c.String(http.StatusNotFound, "Calculation not found")
		return
	}
	// Проверяем, не удалена ли заявка (статус "удалён") – если да, возвращаем 404 или сообщение
	if calc.Status == "удалён" {
		c.String(http.StatusNotFound, "Calculation was deleted")
		return
	}
	c.HTML(http.StatusOK, "calculation.html", gin.H{
		"calculation": calc,
		"items":       items,
	})
}

func (h *Handler) AddToCalculation(c *gin.Context) {
	// Получаем данные из формы
	electrolyteID, _ := strconv.Atoi(c.PostForm("electrolyte_id"))
	volume, _ := strconv.Atoi(c.PostForm("volume"))
	comment := c.PostForm("comment")
	if electrolyteID == 0 || volume == 0 {
		c.String(http.StatusBadRequest, "Missing fields")
		return
	}

	const creatorID = 1 // пока захардкожено

	// Получаем или создаём черновик
	draft, err := h.Repo.GetOrCreateDraft(creatorID)
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}

	// Добавляем раствор
	err = h.Repo.AddElectrolyteToCalculation(draft.ID, uint(electrolyteID), volume, comment)
	if err != nil {
		// Если дубликат, можно проигнорировать или показать сообщение
		c.String(http.StatusConflict, "Этот раствор уже добавлен в заявку")
		return
	}

	// Перенаправляем на страницу заявки
	c.Redirect(http.StatusFound, "/calculation/"+strconv.Itoa(int(draft.ID)))
}

func (h *Handler) DeleteCalculation(c *gin.Context) {
	calcIDStr := c.PostForm("calculation_id")
	calcID, err := strconv.Atoi(calcIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid calculation ID")
		return
	}

	// Логическое удаление через SQL UPDATE (без ORM)
	err = h.Repo.DeleteCalculation(uint(calcID))
	if err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/")
}