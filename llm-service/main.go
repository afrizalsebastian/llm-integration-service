package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/llm-integration-service/llm-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/handlers"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/infrastructure/middleware"
	proto "github.com/afrizalsebastian/llm-integration-service/proto/gen/go/llm/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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

	// start http gateway server
	gatewayServer, err := startHTTPGatewayServer(app)
	if err != nil {
		log.Fatalln("HTTP Gateway Server failed to start")
		os.Exit(1)
	}

	<-signalChan
	log.Println("Shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	shutdownDone := make(chan struct{}, 2)
	// stopping http gateway
	go func() {
		log.Println("shutting down HTTP gateway...")
		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP gateway shutdown error: %v", err)
		}
		shutdownDone <- struct{}{}
	}()

	go func() {
		// stopping grpcServer
		log.Println("Shutting down grpc server")
		grpcServer.GracefulStop()
		shutdownDone <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-shutdownDone
	}

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

func startHTTPGatewayServer(app *bootstrap.Application) (*http.Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := handlers.NewServer(app)
	if err != nil {
		log.Fatal("failed to init server")
		os.Exit(1)
	}

	jsonOption := runtime.WithMarshalerOption("application/json+strict", &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},

		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	errorHandling := runtime.WithErrorHandler(func(ctx context.Context, sm *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "failed",
			"status":  http.StatusInternalServerError,
			"data":    nil,
		})
	})

	grpcMux := runtime.NewServeMux(jsonOption, errorHandling)
	if err := proto.RegisterLlmServiceHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal("Failed to create http gateway")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/llm/", http.StripPrefix("/api/v1/llm", grpcMux))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("404 --- route not found")
		http.Error(w, "404 - route not found", http.StatusNotFound)
	})

	handler := middleware.Chain(
		mux,
		middleware.RecoveryMiddleware(),
		middleware.CORSMiddleware,
		middleware.RequestTracing(),
		middleware.MonitorMiddleware(),
	)

	serve := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.ENV.AppPort),
		Handler: handler,
	}

	// start server
	go func() {
		log.Printf("🚀 Server is running on port %v\n", app.ENV.AppPort)
		if err := serve.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("server failed to start")
			os.Exit(1)
		}
	}()

	return serve, nil
}
