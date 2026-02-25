package ds

type CalculationItem struct {
	ID             uint `gorm:"primaryKey"`
	CalculationID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte"`
	ElectrolyteID  uint `gorm:"not null;uniqueIndex:idx_calc_electrolyte"`
	Volume         int  `gorm:"not null"`      // объём в мл
	Comment        string `gorm:"type:varchar(100)"` // поле "м-м"

	Calculation Calculation `gorm:"foreignKey:CalculationID"`
	Electrolyte Electrolyte  `gorm:"foreignKey:ElectrolyteID"`
}