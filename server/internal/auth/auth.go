package auth

import (
	"context"
	"errors"

	"github.com/Azanul/Next-Watch/graph/model"
)

func GetUserFromContext(ctx context.Context) (*model.User, error) {
	user := ctx.Value("user")

	if user == nil {
		return nil, errors.New("user not found")
	}
	return user.(*model.User), nil
}
