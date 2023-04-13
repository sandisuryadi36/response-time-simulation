package api

import (
	"context"
	"fmt"
	"log"
	"response-time-simulation/server/pb"
	"time"
)

// GRPC Handler
func (s *Server) Store(ctx context.Context, req *pb.StoreRequest) (*pb.StoreResponse, error) {
	result := &pb.StoreResponse{
			Success: true,
			Code:    200,
			Message: "Success storing data",
	}

	log.Println("Storing data...")
	resStartTime := time.Now().Local()

	// create a buffered channel to store the data
	dataCh := make(chan *pb.OrderORM, 1000)

	// start a transaction
	tx := s.provider.BeginTx(ctx)
	defer func() {
			if r := recover(); r != nil {
					tx.Rollback()
					log.Println("Recovered from panic:", r)
			}
	}()

	// process the data in batches using goroutine
	batchSize := 50
	numBatches := (len(req.GetData()) + batchSize - 1) / batchSize
	for i := 0; i < numBatches; i++ {
			start := i * batchSize
			end := (i + 1) * batchSize
			if end > len(req.GetData()) {
					end = len(req.GetData())
			}

			// process the batch in goroutine
			go func(data []*pb.RequestData) {
					for _, val := range data {
						// simulate processing time
						time.Sleep(10 * time.Millisecond)

							timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", val.GetTimestamp(), time.Local)
							if err != nil {
									log.Println("Error parsing timestamp:", err)
									continue
							}
							data := &pb.OrderORM{
									Id:        val.GetId(),
									Cutomer:   val.GetCustomer(),
									Quantity:  val.GetQuantity(),
									Price:     val.GetPrice(),
									RequestId: req.GetRequestId(),
									Timestamp: &timestamp,
							}

							// send the data to the channel
							dataCh <- data
					}
			}(req.GetData()[start:end])
	}

	// process the data in the channel
	for j := 0; j < len(req.GetData()); j++ {
			data := <-dataCh
			if _, err := s.provider.CreateData(ctx, tx, data); err != nil {
					tx.Rollback()
					return nil, err
			}
	}

	// commit the transaction
	if err := tx.Commit().Error; err != nil {
			return nil, err
	}

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
