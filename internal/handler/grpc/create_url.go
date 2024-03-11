package handler

import (
	"context"
	"fmt"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
	pb "github.com/0loff/go_link_shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateShortURL(ctx context.Context, in *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	var response pb.CreateShortURLResponse

	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	if len(in.URL) == 0 {
		logger.Log.Error("The shortening URL cannot be empty")
		return nil, status.Errorf(codes.InvalidArgument, "The shortening URL cannot be empty")
	}

	shortURL, err := h.services.CreateShortURL(ctx, UserID, in.URL)
	response.Result = fmt.Sprintf("%s/%s", h.services.ShortURLHost, shortURL)
	if err != nil {
		logger.Log.Error("The URL for shorten already exists", zap.Error(err))
		return &response, status.Errorf(codes.AlreadyExists, codes.AlreadyExists.String())
	}

	return &response, nil
}
