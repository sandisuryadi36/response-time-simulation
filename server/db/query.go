package db

import (
	"context"
	"response-time-simulation/server/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *GormProvider) CreateData(ctx context.Context, data []*pb.OrderORM) ([]*pb.OrderORM, error) {
	query := p.db_main
	if err := query.CreateInBatches(&data, 100).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error: %v", err)
	}

	return data, nil
}

func (p *GormProvider) ListAllData(ctx context.Context) ([]*pb.OrderORM, error) {
	data := []*pb.OrderORM{}
	query := p.db_main
	if err := query.Order("id").Find(&data).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Error: %v", err)
	}

	return data, nil
}