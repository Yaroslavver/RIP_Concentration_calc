package ds

type Electrolyte struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"type:varchar(100);not null" json:"name"`
	Concentration float64 `gorm:"not null" json:"concentration"`
	Ions         string  `gorm:"type:varchar(50);not null" json:"ions"`
	PH           float64 `json:"ph"`
	Description  string  `gorm:"type:text" json:"description"`
	Image        string  `gorm:"type:varchar(255)" json:"image"`
	Video        string  `gorm:"type:varchar(255)" json:"video"`
	IsDeleted    bool    `gorm:"default:false" json:"-"`
}