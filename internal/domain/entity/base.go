package entity

import "time"

type TimeStamps struct {
	CreatedAt *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at"`
}

// Meta Response
type Meta struct {
	Page          int `json:"page"`
	Limit         int `json:"limit"`
	Total         int `json:"total"`
	TotalFiltered int `json:"total_filtered"`
	LastPage      int `json:"last_page"`
	Draw          int `json:"draw"`
}

type QueryFilter struct {
	Page         int
	Limit        int
	Search       string
	OrderColumn  string
	OrderDir     string
	SearchColumn []string
	Conditions   map[string]interface{} // fleksibel
}
