package handler

import (
	"context"
	"errors"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/repository"
	pb "github.com/0loff/go_link_shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) GetShortURL(ctx context.Context, in *pb.GetShortURLRequest) (*emptypb.Empty, error) {
	URL, err := h.services.GetShortURL(ctx, in.ShortURL)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrURLGone):
			logger.Log.Error("URL has been deleted", zap.Error(err))
			return nil, status.Errorf(codes.Unavailable, "Cannot get URL")

		default:
			logger.Log.Error("Not found", zap.Error(err))
			return nil, status.Errorf(codes.NotFound, "Not found")
		}
	}

	locationHeader := metadata.New(map[string]string{"Location": URL})
	if err := grpc.SetHeader(ctx, locationHeader); err != nil {
		logger.Log.Error("Cannot set response location header", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Cannot set response location header")
	}

	return &emptypb.Empty{}, nil
}
