package interceptors

import (
	"context"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
	"github.com/0loff/go_link_shortener/pkg/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var (
		token string
		uid   string
	)

	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md.Get("token")) > 0 {
		token = md.Get("token")[0]
		userID, err := jwt.GetUserID(token)
		if err == nil {
			uid = userID
		}
	}

	if token == "" || uid == "" {
		ctx = context.WithValue(ctx, utils.ContextKeyUserStatus, "new_user")

		uuid := uuid.New()
		tokenValue, err := jwt.BuildJWTString(uuid)
		if err != nil {
			logger.Log.Error("Cannot create unique auth token", zap.Error(err))
			panic(err)
		}
		token = tokenValue
		uid = uuid.String()
	}

	ctx = context.WithValue(ctx, utils.ContextKeyUserID, uid)

	authHeader := metadata.New(map[string]string{"token": token})
	if err := grpc.SetHeader(ctx, authHeader); err != nil {
		return nil, status.Errorf(codes.Internal, "Cannot send token as auth header")
	}

	return handler(ctx, req)
}
