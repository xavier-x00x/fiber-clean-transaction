package handler

import (
	"crypto/rand"
	"encoding/hex"
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
	"strconv"
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
	}

	usrJwt := dto.UserJwt{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	access, err := jwtutil.GenerateJWT(&usrJwt, 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	refresh, err := jwtutil.GenerateJWT(&usrJwt, 60*24*30*6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   getCookieSecureFlag(),
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

	access, err := jwtutil.GenerateJWT(claims, 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}
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
	userClaims := contextkeys.GetUserC(c.UserContext())
	if userClaims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	user, err := h.usecase.Profile(userClaims.ID)
	if err != nil {
		return ResponseError(c, err)
	}

	usrResponse := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Avatar:    user.Avatar,
		UpdatedAt: *user.UpdatedAt,
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

	userg, err := fetchGoogleUserInfo(req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	randomPass, err := generateRandomPassword(32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate secure password",
		})
	}

	userRequest := &dto.UserRequest{
		Name:     userg.Name,
		Username: userg.Email,
		Email:    userg.Email,
		Password: randomPass,
	}

	if err := h.usecase.Register(userRequest); err != nil {
		return ResponseError(c, err)
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

	userg, err := fetchGoogleUserInfo(req.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.usecase.GoogleProfile(userg.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Email not registered",
		})
	}

	UpdatedAt := *user.UpdatedAt
	if time.Since(UpdatedAt) < time.Hour {
		if userg.Picture != "" {
			h.getPicture(user.ID, userg.ID, userg.Picture)
		}
	}

	usrJwt := dto.UserJwt{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	access, err := jwtutil.GenerateJWT(&usrJwt, 15)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	refresh, err := jwtutil.GenerateJWT(&usrJwt, 60*24*30*6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	cookie := fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		HTTPOnly: true,
		Secure:   getCookieSecureFlag(),
		SameSite: "Lax",
	}

	if req.RememberMe {
		cookie.Expires = time.Now().Add(time.Hour * 24 * 30 * 6)
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"access_token": access,
	})
}

// --- Helper Functions ---

func fetchGoogleUserInfo(accessToken string) (*dto.GoogleUser, error) {
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	reqAPI, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	reqAPI.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(reqAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API error: %s", resp.Status)
	}

	var userg dto.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userg); err != nil {
		return nil, fmt.Errorf("failed to parse user info")
	}

	return &userg, nil
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (h *AuthHandler) getPicture(userID uint, ID string, pic string) {
	respPic, err := http.Get(pic)
	if err != nil {
		return
	}
	defer respPic.Body.Close()

	fileName := fmt.Sprintf("uploads/%s.jpg", ID)
	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer out.Close()

	io.Copy(out, respPic.Body)
	h.usecase.UpdateAvatar(userID, fileName)
}

func getCookieSecureFlag() bool {
	secureFlag := os.Getenv("COOKIE_SECURE")
	if secureFlag == "" {
		return false
	}

	isSecure, err := strconv.ParseBool(secureFlag)
	if err != nil {
		return false
	}

	return isSecure
}
