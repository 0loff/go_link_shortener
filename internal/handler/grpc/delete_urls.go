package handler

import (
	"context"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/utils"
	pb "github.com/0loff/go_link_shortener/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) DeleteShortURLs(ctx context.Context, in *pb.DeleteURLRequest) (*emptypb.Empty, error) {
	var URLsList []string

	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	URLsList = append(URLsList, in.Data...)
	h.services.DelShortURLs(UserID, URLsList)

	return &emptypb.Empty{}, nil
}
