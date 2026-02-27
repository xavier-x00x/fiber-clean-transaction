package entity

type Tax struct {
	Id      uint    `json:"id" gorm:"primarykey"`
	Code    string  `gorm:"size:30;unique;not null"`
	Name    string  `gorm:"size:250;not null"`
	Percent float64 `gorm:"type:decimal(5,2);not null"`
	TaxType string  `gorm:"size:50;not null"` // e.g., "inclusive" or "exclusive"
	Desc    string  `gorm:"type:text"`
	Status  *bool   `gorm:"type:tinyint(1);default:1"`
	TimeStamps
}

func (u *Tax) TableName() string {
	return "taxes"
}
