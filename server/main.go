package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"response-time-simulation/server/api"
	"response-time-simulation/server/pb"
)

func main() {
	// migrate DB
	migrateDB()

	// start DB connection
	startDBConnection()

	// Inisialisasi gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingMiddleware),
	)

	apiServ := api.New(
		db_main,
	)
	// Registrasi handler ke gRPC server
	pb.RegisterApiServiceServer(grpcServer, apiServ)

	// Registrasi reflection service untuk debugging
	reflection.Register(grpcServer)

	// Inisialisasi listener untuk HTTP gateway
	httpListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Inisialisasi gRPC-gateway Mux
	gwMux := runtime.NewServeMux()

	// Register HTTP handler for gRPC service
	err = pb.RegisterApiServiceHandlerFromEndpoint(context.Background(), gwMux, fmt.Sprintf("localhost:%d", 9090), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("Failed to register HTTP gateway: %v", err)
	}

	// Inisialisasi HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: loggingHTTPMiddleware(gwMux),
	}

	// Mulai server gRPC dan HTTP API
	go func() {
		log.Printf("Starting gRPC server on localhost:9090")
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	go func() {
		log.Printf("Starting HTTP server on localhost:8080")
		err = httpServer.Serve(httpListener)
		if err != nil {
			log.Fatalf("Failed to serve HTTP server: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch

	closeDBMain()
}

func GetEnv(key, fallback string) string {
	godotenv.Load(".env")
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

