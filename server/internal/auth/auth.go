package auth

import (
	"context"
	"errors"

	"github.com/Azanul/Next-Watch/internal/models"
)

func GetUserFromContext(ctx context.Context) (*models.User, error) {
	user := ctx.Value("user")

	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.(*models.User), nil
}
