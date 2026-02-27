package entity

type Store struct {
	Id      uint    `json:"id"`
	Code    string  `gorm:"size:30;unique;not null"`
	Name    string  `gorm:"size:250;not null"`
	Npwp    string  `gorm:"size:250"`
	Address string  `gorm:"type:text"`
	Phone   string  `gorm:"size:30"`
	Email   string  `gorm:"size:250"`
	Phone2  string  `gorm:"size:30"`
	Email2  string  `gorm:"size:250"`
	Disc    float64 `gorm:"type:decimal(10,2);default:0"`
	Status  bool    `gorm:"type:tinyint(1);default:1"`
	TimeStamps
}