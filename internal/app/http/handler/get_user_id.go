package handler

import (
	"context"
	"github.com/aleksey-kombainov/url-shortener.git/internal/app/user"
)

func getUserIDFromCtx(ctx context.Context) string {
	if val := ctx.Value(user.CtxUserIDKey); val != nil {
		return val.(string)
	}
	return ""
}
