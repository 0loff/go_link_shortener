// Утилита для работы с ключаим параметров, передаваемых с помощью контекста
package utils

import "context"

type contextKey string

// Ключи параметров передаваемых с помощью контекста
var (
	ContextKeyUserID     = contextKey("uid")
	ContextKeyUserStatus = contextKey("user status")
)

// Вывод ключа параметра передаваемого с помощью контекста в формате String
func (c contextKey) String() string {
	return string(c)
}

// Получение ключа ID пользователя
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(ContextKeyUserID).(string)
	return UserID, ok
}

// Получение ключа статуса пользователя
func GetUserStatusFromContext(ctx context.Context) (string, bool) {
	UserStatus, ok := ctx.Value(ContextKeyUserStatus).(string)
	return UserStatus, ok
}
