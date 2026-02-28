package dto

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type GoogleAuthRequest struct {
	AccessToken string `json:"access_token"`
	RememberMe  bool   `json:"remember_me"`
}
