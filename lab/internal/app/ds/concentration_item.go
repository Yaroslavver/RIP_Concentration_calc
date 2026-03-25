package ds

type ConcentrationItem struct {
	ID             uint `gorm:"primaryKey" json:"id"`
	ConcentrationID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte" json:"concentration_id"`
	ElectrolyteID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte" json:"electrolyte_id"`
	Volume         int  `gorm:"not null" json:"volume"`      // объём в мл, поле "м-м"
	Comment        string `gorm:"type:varchar(100)" json:"comment"` // 

	Concentration Concentration `gorm:"foreignKey:ConcentrationID" json:"-"`
	Electrolyte Electrolyte  `gorm:"foreignKey:ElectrolyteID" json:"electrolyte"`
}