package entity

type Category struct {
	ID          uint   `gorm:"primarykey"`
	StoreCode   string `gorm:"size:30;not null"`
	Code        string `gorm:"size:30;unique;not null"`
	Name        string `gorm:"size:250;not null"`
	Description string `gorm:"type:text"`
	Status      *bool  `gorm:"type:tinyint(1);default:1"`
	TimeStamps
}
