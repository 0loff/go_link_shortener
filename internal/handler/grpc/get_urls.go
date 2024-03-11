package handler

import (
	"context"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
	pb "github.com/0loff/go_link_shortener/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetShortURLs(ctx context.Context, in *empty.Empty) (*pb.GetShortURLsResponse, error) {
	var response pb.GetShortURLsResponse
	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	UserStatus, ok := utils.GetUserStatusFromContext(ctx)
	if ok && UserStatus == "new_user" {
		return nil, status.Errorf(codes.Unauthenticated, "User is unauthorized")
	}

	resp := h.services.GetShortURLs(ctx, UserID)
	if len(resp) == 0 {
		return nil, status.Errorf(codes.NotFound, "Not found URLs by user")
	}

	for _, relation := range resp {
		response.Data = append(response.Data, &pb.Relation{
			ShortUrl:    relation.ShortURL,
			OriginalUrl: relation.OriginalURL,
		})
	}

	return &response, nil
}
