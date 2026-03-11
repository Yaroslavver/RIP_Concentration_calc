package repository

import (
	"context"
	"fmt"
	"lab/internal/app/ds"
	"mime/multipart"
	"net/http"
	//"strconv"

	"github.com/minio/minio-go/v7"
	//"github.com/sirupsen/logrus"
)

func (r *Repository) GetAllElectrolytes() ([]ds.Electrolyte, error) {
	var electrolytes []ds.Electrolyte
	err := r.db.Where("is_deleted = ?", false).Find(&electrolytes).Error
	return electrolytes, err
}

func (r *Repository) SearchElectrolytesByName(name string) ([]ds.Electrolyte, error) {
	var electrolytes []ds.Electrolyte
	err := r.db.Where("name ILIKE ? AND is_deleted = ?", "%"+name+"%", false).Find(&electrolytes).Error
	return electrolytes, err
}

func (r *Repository) GetElectrolyteByID(id uint) (*ds.Electrolyte, error) {
	var electrolyte ds.Electrolyte
	err := r.db.Where("id = ? AND is_deleted = ?", id, false).First(&electrolyte).Error
	if err != nil {
		return nil, err
	}
	return &electrolyte, nil
}

func (r *Repository) CreateElectrolyte(e *ds.Electrolyte) error {
	// Убеждаемся, что ID не передан (GORM сам сгенерирует)
	e.ID = 0
	return r.db.Create(e).Error
}

func (r *Repository) UpdateElectrolyte(e *ds.Electrolyte) error {
	return r.db.Model(e).Updates(map[string]interface{}{
		"name":         e.Name,
		"concentration": e.Concentration,
		"ions":         e.Ions,
		"ph":           e.PH,
		"description":  e.Description,
		"image":        e.Image,
		"video":        e.Video,
	}).Error
}

func (r *Repository) DeleteElectrolyte(id uint) error {
	return r.db.Model(&ds.Electrolyte{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UploadElectrolyteFile загружает файл в MinIO и обновляет ссылку в БД
func (r *Repository) UploadElectrolyteFile(electrolyteID uint, fileHeader *multipart.FileHeader, fileType string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Определяем Content-Type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	file.Seek(0, 0)

	// Генерируем имя файла
	objectName := fmt.Sprintf("electrolyte_%d_%s", electrolyteID, fileType)
	switch contentType {
	case "image/jpeg":
		objectName += ".jpg"
	case "image/png":
		objectName += ".png"
	case "video/mp4":
		objectName += ".mp4"
	default:
		objectName += ".bin"
	}

	// Загружаем в MinIO с контекстом
	ctx := context.Background()
	_, err = r.minio.PutObject(ctx, r.minioBucket, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// Формируем URL
	scheme := "http"
	if r.minioUseSSL {
		scheme = "https"
	}
	url := fmt.Sprintf("%s://%s/%s/%s", scheme, r.minioEndpoint, r.minioBucket, objectName)

	// Обновляем поле в БД
	var updateField string
	if fileType == "image" {
		updateField = "image"
	} else {
		updateField = "video"
	}
	err = r.db.Model(&ds.Electrolyte{}).Where("id = ?", electrolyteID).Update(updateField, url).Error
	if err != nil {
		// При ошибке удаляем файл из MinIO (с контекстом)
		r.minio.RemoveObject(ctx, r.minioBucket, objectName, minio.RemoveObjectOptions{})
		return "", err
	}
	return url, nil
}