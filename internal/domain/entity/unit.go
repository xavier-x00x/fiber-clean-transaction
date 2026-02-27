package entity

type Unit struct {
	Id     uint   `json:"id" gorm:"primarykey"`
	Code   string `gorm:"size:30;unique;not null"`
	Name   string `gorm:"size:250;not null"`
	Status *bool  `gorm:"type:tinyint(1);default:1"`
	TimeStamps
}

func (u *Unit) TableName() string {
	return "measurement_units"
}