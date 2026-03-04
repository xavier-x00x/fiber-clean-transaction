package entity

type NumberSequence struct {
	ID         uint   `json:"id" gorm:"primarykey"`
	Prefix     string `gorm:"size:10;not null;uniqueIndex:uq_prefix_period"`
	Period     string `gorm:"size:4;not null;uniqueIndex:uq_prefix_period"` // format: YYMM
	LastNumber int    `gorm:"default:0"`
	TimeStamps
}

func (n *NumberSequence) TableName() string {
	return "number_sequences"
}
