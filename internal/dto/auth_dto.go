package dto

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type VerifyUserRequest struct {
	Token string `json:"token"`
}

type GoogleAuthRequest struct {
	AccessToken string `json:"access_token"`
	RememberMe  bool   `json:"remember_me"`
}

type UserRequestX struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}