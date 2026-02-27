package middleware

import (
	"context"
	"fiber-clean-transaction/internal/contextkeys"
	"fiber-clean-transaction/pkg/jwtutil"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware() fiber.Handler {
	return AuthMiddleware
}

func AuthMiddleware(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error":  "missing or invalid authorization header",
		})
	}

	// Ambil token setelah "Bearer "
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	// println(tokenStr)

	claims, err := jwtutil.ValidateJWT(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error":  err.Error(),
		})
	}

	// Simpan claims ke context biar bisa diakses di handler
	// mulai dari context fiber
	ctx := c.UserContext()

	ctx = context.WithValue(ctx, contextkeys.KeyUser, claims)

	c.SetUserContext(ctx)

	return c.Next()
}
