package api

import (

	"response-time-simulation/server/db"
	"response-time-simulation/server/pb"

	"gorm.io/gorm"
)

// Setup Server
type Server struct {
	provider *db.GormProvider
	pb.ApiServiceServer
}

func New(db01 *gorm.DB) *Server {
	return &Server{
		provider: db.NewProvider(db01),
		ApiServiceServer: nil,
	}
}
