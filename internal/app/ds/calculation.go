package ds

import (
	"time"
)

type Calculation struct {
	ID          uint      `gorm:"primaryKey"`
	Status      string    `gorm:"type:varchar(20);not null;default:'черновик'"` // черновик, удалён, сформирован, завершён, отклонён
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time
	FinishedAt  *time.Time
	CreatorID   uint      `gorm:"not null"`
	ModeratorID *uint
	Result      string    `gorm:"type:varchar(100)"` // например "[H⁺] = 0.045 моль/л, pH = 1.35"
}