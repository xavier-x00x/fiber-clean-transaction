package handler

import (
	"encoding/json"
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/internal/dto"
	"fiber-clean-transaction/internal/usecase"
	"fiber-clean-transaction/pkg/jwtutil"
	"fiber-clean-transaction/pkg/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	usecase *usecase.UserUsecase
}

func NewAuthHandler(u *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	userRequest := new(dto.UserRequest)

	if err := c.BodyParser(&userRequest); err != nil {
		error := utils.BadRequest(err.Error())
		return ResponseError(c, error)
	}

	if err := h.usecase.Register(userRequest); err != nil {
		return ResponseError(c, err)
	}

	return c.JSON(fiber.Map{"message": "user registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	loginInput := new(dto.LoginInput)

	if err := c.BodyParser(&loginInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if loginInput.Email == "" || loginInput.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing email or password",
		})
	}

	user, err := h.usecase.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
		// return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	usrJwt := dto.UserJwt{
		ID:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	// Access token
	access, _ := jwtutil.GenerateJWT(&usrJwt, 15) // 15 minutes
	// Refresh token
	refresh, _ := jwtutil.GenerateJWT(&usrJwt, 60*24*30*6) // 6 months

	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}

	if loginInput.Remember {
		cookie.Expires = time.Now().Add(time.Hour * 24 * 30 * 6)
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"access_token": access,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing refresh token",
		})
	}

	claims, err := jwtutil.ValidateJWT(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid refresh token",
		})
	}

	access, _ := jwtutil.GenerateJWT(claims, 15) // 15 minutes
	return c.JSON(fiber.Map{
		"access_token": access,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("refresh_token")
	return c.JSON(fiber.Map{
		"message": "Logout success",
	})
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	// get context user id
	userClaims := contextkeys.GetUser(c)
	user_id := userClaims.ID

	println(user_id)

	user, _ := h.usecase.Profile(user_id)

	usrResponse := dto.UserResponse{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		UpdatedAt: user.UpdatedAt,
	}

	return c.JSON(usrResponse)
}

func (h *AuthHandler) GoogleRegister(c *fiber.Ctx) error {
	req := new(dto.GoogleAuthRequest)

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing access_token",
		})
	}

	// Panggil Google Userinfo API
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	reqAPI, _ := http.NewRequest("GET", userInfoURL, nil)
	reqAPI.Header.Set("Authorization", "Bearer "+req.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(reqAPI)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user info",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": fmt.Sprintf("Google API error: %s", resp.Status),
		})
	}

	var userg dto.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user info",
		})
	}

	userRequest := &dto.UserRequest{
		Name:     userg.Name,
		Username: userg.Email,
		Email:    userg.Email,
		Password: "password",
	}

	if err := h.usecase.Register(userRequest); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user registered"})
}

func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	var req dto.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.AccessToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing access_token",
		})
	}

	// Panggil Google Userinfo API
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	reqAPI, _ := http.NewRequest("GET", userInfoURL, nil)
	reqAPI.Header.Set("Authorization", "Bearer "+req.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(reqAPI)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user info",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": fmt.Sprintf("Google API error: %s", resp.Status),
		})
	}

	var userg dto.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user info",
		})
	}

	user, err := h.usecase.GoogleProfile(userg.Email)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email not registered",
		})
	}

	UpdatedAt := time.Time(user.UpdatedAt)
	// Jika user terakhir update lebih dari 1 jam sebelumnya, ambil gambar
	if time.Since(UpdatedAt) < time.Hour {
		// ambil gambar
		if userg.Picture != "" {
			h.getPicture(user.Id, userg.ID, userg.Picture)
		}
	}

	usrJwt := dto.UserJwt{
		ID:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	// Access token
	access, _ := jwtutil.GenerateJWT(&usrJwt, 15) // 15 minutes
	// Refresh token
	refresh, _ := jwtutil.GenerateJWT(&usrJwt, 60*24*30*6) // 6 months

	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}

	if req.RememberMe {
		cookie.Expires = time.Now().Add(time.Hour * 24 * 30 * 6) // 6 months
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"access_token": access,
	})
}

func (h *AuthHandler) getPicture(userId uint, id string, pic string) (string, error) {
	var avatar string
	respPic, err := http.Get(pic)
	if err == nil {
		defer respPic.Body.Close()
		fileName := fmt.Sprintf("uploads/%s.jpg", id)
		out, _ := os.Create(fileName)
		defer out.Close()

		io.Copy(out, respPic.Body)
		avatar = fileName // ganti dengan path lokal
		h.usecase.UpdateAvatar(userId, avatar)
	}
	return string(avatar), err
}
