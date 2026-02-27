package dto

import "time"

type GetTimeStamps struct {
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type MetaRequest struct {
	Page  		int    `json:"page"`
	Limit 		int    `json:"limit"`
	Search 		string `json:"search"`
	OrderColumn string `json:"order_column"`
	OrderDir 	string `json:"order_dir"`
}

type JsonResponse struct {
	Success bool         `json:"success"`
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data"`
	Meta    *interface{} `json:"meta"`
	Error   string 	 `json:"error"`
	Errors  *interface{} `json:"errors"`
}

type MsgError struct {
	Status  int  `json:"status"`
	Error string `json:"error"`
}