package contextkeys

import (
	"context"
	"fiber-clean-transaction/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type contextKey string

const (
	KeyID    contextKey = "user_id"
	KeyEmail contextKey = "email"
	KeyRole  contextKey = "role"
	KeyStore contextKey = "store"
)

func GetUserCtx(ctx context.Context) *dto.UserJwt {
	userClaims := dto.UserJwt{
		ID:    ctx.Value(KeyID).(uint),
		Email: ctx.Value(KeyEmail).(string),
		Role:  ctx.Value(KeyRole).(string),
		Store: ctx.Value(KeyStore).(string),
	}
	return &userClaims
}

func GetUser(c *fiber.Ctx) *dto.UserJwt {
	ctx := c.UserContext()
	// userClaims := dto.UserJwt{
	// 	ID:    ctx.Value(KeyID).(uint),
	// 	Email: ctx.Value(KeyEmail).(string),
	// 	Role:  ctx.Value(KeyRole).(string),
	// }
	userClaims, _ := ctx.Value(KeyUser).(*dto.UserJwt)
	return userClaims
}

type userKey struct{}

var KeyUser = userKey{}

func GetUserC(ctx context.Context) *dto.UserJwt {
	user, _ := ctx.Value(KeyUser).(*dto.UserJwt)
	return user
}
