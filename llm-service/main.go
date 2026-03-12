package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/llm-integration-service/llm-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/handlers"
	proto "github.com/afrizalsebastian/llm-integration-service/proto/gen/go/llm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := bootstrap.NewApp()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start grpc server
	grpcServer, err := startGRPCServer(app)
	if err != nil {
		log.Fatalln("GRPC Server failed to start")
		os.Exit(1)
	}

	_ = <-signalChan
	log.Println("Shutting down server")

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// stopping grpcServer
	grpcServer.GracefulStop()

	log.Println("✅ Server exited gracefully")
}

func startGRPCServer(app *bootstrap.Application) (*grpc.Server, error) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, err := net.Listen("tcp", ":"+app.ENV.GrpcPort)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer()
	grpcHandlers, err := handlers.NewServer(app)
	if err != nil {
		log.Fatal("failed to init server")
		os.Exit(1)
	}
	proto.RegisterLlmServiceServer(server, grpcHandlers)
	reflection.Register(server)

	go func() {
		log.Printf("🚀 GRPC Server is running on port %s\n", app.ENV.GrpcPort)
		if err := server.Serve(lis); err != nil {
			log.Fatalln("GRPC Server failed to start")
			os.Exit(1)
		}
	}()

	return server, nil
}
