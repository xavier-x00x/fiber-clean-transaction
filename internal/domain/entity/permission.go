package entity

type Permission struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Path string `gorm:"unique" json:"path"`
	Name string `json:"name"`
	TimeStamps
}
