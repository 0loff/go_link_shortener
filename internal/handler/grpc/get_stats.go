package handler

import (
	"context"
	"net"

	"github.com/0loff/go_link_shortener/internal/logger"
	pb "github.com/0loff/go_link_shortener/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetStats(ctx context.Context, in *empty.Empty) (*pb.GetStatsResponse, error) {
	_, ts, err := net.ParseCIDR(h.trustedSubnet)
	if err != nil {
		logger.Log.Error("The value of the trusted subnet could not be parsed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error during parse CIDR value")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("X-Real-IP")) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "ip is not in trusted subnet")
	}

	ip := net.ParseIP(md.Get("X-Real-IP")[0])

	if !ts.Contains(ip) {
		return nil, status.Errorf(codes.PermissionDenied, "ip is not in trusted subnet")
	}

	metrics := h.services.GetStatistics()

	return &pb.GetStatsResponse{
		Urls:  int32(metrics.Urls),
		Users: int32(metrics.Users),
	}, nil
}
