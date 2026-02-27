package dto

import (
	"encoding/json"
	"fiber-clean-transaction/pkg/utils"
	"fmt"
	"os"
	"time"
)

type UserRequest struct {
	Name      string    `validate:"required,min=3" json:"name"`
	Email     string    `validate:"required,email,unique=users:email" json:"email"`
	Username  string    `validate:"required,unique=users:username" json:"username"`
	Password  string    `validate:"required" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	Avatar    string         `json:"avatar"`
	UpdatedAt utils.JSONTime `json:"updated_at"`
}

func (UserResponse) TableName() string {
	return "users"
}

// Custom JSON Marshal
func (u UserResponse) MarshalJSON() ([]byte, error) {
	type Alias UserResponse
	return json.Marshal(&struct {
		*Alias
		Avatar string `json:"avatar"`
		// UpdatedAt string `json:"updated_at"`
	}{
		// UpdatedAt: func() string {
		// 	if !u.UpdatedAt.Valid {
		// 		return "-"
		// 	}
		// 	return u.UpdatedAt.Time.Format("02/01/2006 15:04:05")
		// }(),
		Avatar: func() string {
			if u.Avatar != "" {
				return fmt.Sprintf("%s%s", os.Getenv("BASE_URL"), u.Avatar)
			}
			return u.Avatar
		}(),
		Alias: (*Alias)(&u),
	})
}

type UserJwt struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Store string `json:"store"`
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email" validate:"unique=users:email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
