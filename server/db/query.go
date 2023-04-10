package db

import (
	"context"
	"response-time-simulation/server/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *GormProvider) CreateData(ctx context.Context, data *pb.Order) (*pb.OrderORM, error) {
	dataORM, _ := data.ToORM(ctx)
	query := p.db_main.Debug()
	if err := query.Create(&dataORM).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error: %v", err)
	}

	return &dataORM, nil
}

func (p *GormProvider) ListAllData(ctx context.Context) ([]*pb.OrderORM, error) {
	data := []*pb.OrderORM{}
	query := p.db_main.Debug()
	if err := query.Find(&data).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error: %v", err)
	}

	return data, nil
}