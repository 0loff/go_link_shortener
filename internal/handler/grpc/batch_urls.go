package handler

import (
	"context"

	"github.com/0loff/go_link_shortener/internal/logger"
	"github.com/0loff/go_link_shortener/internal/models"
	"github.com/0loff/go_link_shortener/internal/utils"
	pb "github.com/0loff/go_link_shortener/proto"
)

// BatchShortURLs - Handler for batch insertion of short URLs by user
func (h *Handler) BatchShortURLs(ctx context.Context, in *pb.CreateBatchShortURLsRequest) (*pb.CreateBatchShortURLsResponse, error) {
	var response pb.CreateBatchShortURLsResponse
	entries := []models.BatchURLRequestEntry{}

	UserID, ok := utils.GetUserIDFromContext(ctx)
	if !ok {
		logger.Log.Error("Cannot get UserID from context")
	}

	for _, entry := range in.Data {
		entries = append(entries, models.BatchURLRequestEntry{
			CorrelationID: entry.CorrelationId,
			OriginalURL:   entry.OriginalUrl,
		})
	}

	resp := h.services.SetBatchShortURLs(ctx, UserID, entries)

	for _, entry := range resp {
		response.Data = append(response.Data, &pb.CorrelationShortURL{
			CorrelationId: entry.CorrelationID,
			ShortUrl:      entry.ShortURL,
		})
	}

	return &response, nil
}
