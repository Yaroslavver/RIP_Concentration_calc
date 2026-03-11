package ds

type ConcentrationItem struct {
	ID             uint `gorm:"primaryKey"`
	ConcentrationID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte"`
	ElectrolyteID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte"`
	Volume         int  `gorm:"not null"`      // объём в мл
	Comment        string `gorm:"type:varchar(100)"` // поле "м-м"

	Concentration Concentration `gorm:"foreignKey:ConcentrationID"`
	Electrolyte Electrolyte  `gorm:"foreignKey:ElectrolyteID"`
}