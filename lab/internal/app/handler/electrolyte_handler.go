package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab/internal/app/ds"
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
	c.JSON(http.StatusOK, gin.H{"data": electrolytes})
}

func (h *Handler) GetElectrolyte(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	electrolyte, err := h.Repo.GetElectrolyteByID(uint(id))
	if err != nil {
		h.errorJSON(c, http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": electrolyte})
}

func (h *Handler) CreateElectrolyte(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		h.errorJSON(c, http.StatusBadRequest, err)
		return
	}

	// Получаем значения из формы
	name := c.Request.FormValue("name")
	concentrationStr := c.Request.FormValue("concentration")
	ions := c.Request.FormValue("ions")
	phStr := c.Request.FormValue("ph")
	description := c.Request.FormValue("description")

	// Валидация обязательных полей
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if concentrationStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "concentration is required"})
		return
	}
	if ions == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ions is required"})
		return
	}

	// Парсинг числовых значений
	concentration, err := strconv.ParseFloat(concentrationStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid concentration"})
		return
	}

	ph := 0.0
	if phStr != "" {
		ph, err = strconv.ParseFloat(phStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ph"})
			return
		}
	}

	// Создаём объект электролита (без ID - он сгенерируется автоматически)
	electrolyte := ds.Electrolyte{
		Name:         name,
		Concentration: concentration,
		Ions:         ions,
		PH:           ph,
		Description:  description,
		Image:        "", // временно пусто
		Video:        "", // временно пусто
	}

	// Сохраняем в БД (ID сгенерируется автоматически)
	if err := h.Repo.CreateElectrolyte(&electrolyte); err != nil {
		h.errorJSON(c, http.StatusInternalServerError, err)
		return
	}

	// Загружаем файлы, если есть
	imageFile, err := c.FormFile("image")
	if err == nil {
		imageURL, err := h.Repo.UploadElectrolyteFile(electrolyte.ID, imageFile, "image")
		if err != nil {
			logrus.Error("failed to upload image:", err)
		} else {
			electrolyte.Image = imageURL
		}
	}

	videoFile, err := c.FormFile("video")
	if err == nil {
		videoURL, err := h.Repo.UploadElectrolyteFile(electrolyte.ID, videoFile, "video")
		if err != nil {
			logrus.Error("failed to upload video:", err)
		} else {
			electrolyte.Video = videoURL
		}
	}

	// Обновляем запись с URL файлов
	if electrolyte.Image != "" || electrolyte.Video != "" {
		h.Repo.UpdateElectrolyte(&electrolyte)
	}

	// Получаем обновлённые данные
	updated, _ := h.Repo.GetElectrolyteByID(electrolyte.ID)
	c.JSON(http.StatusCreated, gin.H{"data": updated})
}