package dto

type NuxtMenuResponse struct {
	Type   string `json:"type"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Icon   string `json:"icon"`
	Parent string `json:"parent"`
	To     string `json:"to"`
}
