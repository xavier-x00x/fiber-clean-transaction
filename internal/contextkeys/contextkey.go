package contextkeys

import (
	"context"
	"fiber-clean-transaction/internal/dto"
)

type userKey struct{}

var KeyUser = userKey{}

func GetUserC(ctx context.Context) *dto.UserJwt {
	user, _ := ctx.Value(KeyUser).(*dto.UserJwt)
	return user
}
