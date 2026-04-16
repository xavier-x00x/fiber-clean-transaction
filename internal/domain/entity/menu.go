package entity

type NuxtMenu struct {
	ID     int    `json:"id" gorm:"primarykey"`
	Type   string `json:"type" gorm:"size:30;not null"`
	Code   string `json:"code" gorm:"size:30;not null"`
	Title  string `json:"title" gorm:"size:250;not null"`
	Icon   string `json:"icon" gorm:"size:100"`
	Parent string `json:"parent" gorm:"size:30"`
	To     string `json:"to" gorm:"size:250"`
	Active int    `json:"active" gorm:"type:tinyint(1);default:1"`
}

func (NuxtMenu) TableName() string {
	return "nuxt_menu"
}

// type NuxtMenuResponse struct {
// 	Type   string `json:"type"`
// 	Code   string `json:"code"`
// 	Title  string `json:"title"`
// 	Icon   string `json:"icon"`
// 	Parent string `json:"parent"`
// 	To     string `json:"to"`
// }

// func (NuxtMenuResponse) TableName() string {
// 	return "nuxt_menu"
// }
