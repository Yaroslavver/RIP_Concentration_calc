package ds

type Electrolyte struct {
	ID           uint    `gorm:"primaryKey"`
	Name         string  `gorm:"type:varchar(100);not null"`
	Concentration float64 `gorm:"not null"`
	Ions         string  `gorm:"type:varchar(50);not null"`
	PH           float64
	Description  string  `gorm:"type:text"`
	Image        string  `gorm:"type:varchar(255)"`
	Video        string  `gorm:"type:varchar(255)"`
	IsDeleted    bool    `gorm:"default:false"`
}