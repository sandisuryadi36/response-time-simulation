package api

import (
	"context"
	"log"
	"response-time-simulation/server/pb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPC Handler
func (s *Server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	result := &pb.StoreResponse{
		Success: true,
		Code:    200,
		Message: "Success storing data",
	}
	
	log.Println("Storing data...")
	for _, val := range req.GetData() {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", val.GetTimestamp(), time.Local)
		if err != nil {
			log.Println("Error parsing timestamp:", err)
		}
		data := &pb.Order{
			Id: val.GetId(),
			Cutomer: val.GetCustomer(),
			Quantity: val.GetQuantity(),
			Price: val.GetPrice(),
			RequestId: req.GetRequestId(),
			Timestamp: timestamppb.New(timestamp),
		}
		res, err := s.provider.CreateData(ctx, data)
		if err != nil {
			return nil, err
		}
		log.Println("Data Stored with ID:", res.Id)
	}

	return result, nil
}

func (s *Server) List(ctx context.Context, req *pb.Empty) (*pb.ListResponse, error) {
	result := &pb.ListResponse{
		Success: true,
		Code:    200,
		Message: "List all data",
	}

	dataORM, err := s.provider.ListAllData(ctx)
	if err != nil {
		return nil, err
	}
	res := []*pb.ListData{}
	for _, val := range dataORM {
		dataPb, _ := val.ToPB(ctx)
		res = append(res, &pb.ListData{
			Id: dataPb.Id,
			Cutomer: dataPb.Cutomer,
			Quantity: dataPb.Quantity,
			Price: dataPb.Price,
			RequestId: dataPb.RequestId,
			Timestamp: dataPb.Timestamp,
		})
	}
	result.Data = res
	return result, nil
}
