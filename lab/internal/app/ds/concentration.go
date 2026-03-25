package ds

import (
	"time"
)

type Concentration struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Status      string    `gorm:"type:varchar(20);not null;default:'черновик'" json:"status"` // черновик, удалён, сформирован, завершён, отклонён
	CreatedAt   time.Time `gorm:"not null" json:"-"` //json:"created_at"`
	UpdatedAt   time.Time `json:"-"` //`json:"updated_at"`
	FinishedAt  *time.Time `json:"-"` //`json:"finished_at,omitempty"`
	CreatorID   uint      `gorm:"not null" json:"-"`//json:"creator_id"`
	ModeratorID *uint  `json:"-"`//`json:"moderator_id,omitempty"`
	Result      string    `gorm:"type:varchar(100)" json:"result"` // например "[H⁺] = 0.045 моль/л, pH = 1.35"
	Description string     `gorm:"type:text" json:"description"` // новое поле

	Creator   User `gorm:"foreignKey:CreatorID" json:"-"` // было json:"creator"`
	Moderator User `gorm:"foreignKey:ModeratorID" json:"-"` // было json:"moderator,omitempty"`
}