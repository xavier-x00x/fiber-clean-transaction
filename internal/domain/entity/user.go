package entity

import "fiber-clean-transaction/pkg/utils"

type User struct {
	Id        uint           `json:"id"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Role      string         `json:"role"`
	Avatar    string         `json:"avatar"`
	CreatedAt utils.JSONTime `gorm:"<-:create;autoCreateTime" json:"created_at"`
	UpdatedAt utils.JSONTime `gorm:"autoUpdateTime" json:"updated_at"`
}
