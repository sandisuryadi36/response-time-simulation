package api

import (
	"context"
	"fmt"
	"log"
	"response-time-simulation/server/pb"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

// GRPC Handler
func (s *Server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	result := &pb.StoreResponse{
		Success: true,
		Code:    200,
		Message: "Success storing data",
	}

	var wg sync.WaitGroup

	log.Println("Storing data...")
	resStartTime := time.Now().Local()
	dataToStore := []*pb.OrderORM{}
	for _, val := range req.GetData() {
		goData := proto.Clone(val).(*pb.RequestData)

		wg.Add(1)
		go func(val *pb.RequestData) {
			defer wg.Done()
			startT := time.Now().Local()

			// time process simulation
			time.Sleep(time.Duration(10) * time.Millisecond)

			timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", val.GetTimestamp(), time.Local)
			if err != nil {
				log.Println("Error parsing timestamp:", err)
			}
			data := &pb.OrderORM{
				Id:        val.GetId(),
				Cutomer:   val.GetCustomer(),
				Quantity:  val.GetQuantity(),
				Price:     val.GetPrice(),
				RequestId: req.GetRequestId(),
				Timestamp: &timestamp,
			}

			dataToStore = append(dataToStore, data)
			processT := time.Now().Local().Sub(startT)
			log.Printf("Data ID: %v process time done in %v", data.Id,processT)
		}(goData)
	}
	wg.Wait()
	s.provider.CreateData(ctx, dataToStore)

	// wg.Wait()
	resTime := time.Now().Local().Sub(resStartTime)
	log.Println("Success store data in", resTime)
	result.Message = fmt.Sprintf("Success store data in %v", resTime)

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
	result.DataCount = 0
	for _, val := range dataORM {
		dataPb, _ := val.ToPB(ctx)
		res = append(res, &pb.ListData{
			Id:        dataPb.Id,
			Cutomer:   dataPb.Cutomer,
			Quantity:  dataPb.Quantity,
			Price:     dataPb.Price,
			RequestId: dataPb.RequestId,
			Timestamp: dataPb.Timestamp,
		})
		result.DataCount++
	}
	result.Data = res
	return result, nil
}
