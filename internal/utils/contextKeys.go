package utils

import "context"

type contextKey string

var (
	ContextKeyUserID     = contextKey("uid")
	ContextKeyUserStatus = contextKey("user status")
)

func (c contextKey) String() string {
	return string(c)
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(ContextKeyUserID).(string)
	return UserID, ok
}

func GetUserStatusFromContext(ctx context.Context) (string, bool) {
	UserStatus, ok := ctx.Value(ContextKeyUserStatus).(string)
	return UserStatus, ok
}
