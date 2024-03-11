package handler

import (
	"github.com/0loff/go_link_shortener/internal/service"

	pb "github.com/0loff/go_link_shortener/proto"
)

type Handler struct {
	pb.UnimplementedShrotenerServer

	services      *service.Service
	trustedSubnet string
}

func NewHandler(s *service.Service, ts string) *Handler {
	return &Handler{
		services:      s,
		trustedSubnet: ts,
	}
}
