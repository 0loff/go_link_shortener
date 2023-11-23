package utils

import "context"

type contextKey string

var (
	ContextKeyUserID contextKey
)

func (c contextKey) String() string {
	return string(c)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(ContextKeyUserID).(string)
	return UserID, ok
}
