package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/handlers"
	"github.com/afrizalsebastian/llm-integration-service/llm-service/infrastructure/middleware"
	proto "github.com/afrizalsebastian/llm-integration-service/proto/gen/go/llm/v1"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	app := bootstrap.NewApp()

	l := logger.New()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start grpc server
	grpcServer, err := startGRPCServer(app)
	if err != nil {
		l.Warn("GRPC Server failed to start").Msg()
		os.Exit(1)
	}

	// start http gateway server
	gatewayServer, err := startHTTPGatewayServer(app)
	if err != nil {
		l.Warn("HTTP Gateway Server failed to start").Msg()
		os.Exit(1)
	}

	<-signalChan
	l.Info("Shutting down server").Msg()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	shutdownDone := make(chan struct{}, 2)
	// stopping http gateway
	go func() {
		l.Info("shutting down HTTP gateway...").Msg()
		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
			l.Errorf("HTTP gateway shutdown error: %v", err).Msg()
		}
		shutdownDone <- struct{}{}
	}()

	go func() {
		// stopping grpcServer
		l.Info("Shutting down grpc server").Msg()
		grpcServer.GracefulStop()
		shutdownDone <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		<-shutdownDone
	}

	l.Info("✅ Server exited gracefully").Msg()
}

func startGRPCServer(app *bootstrap.Application) (*grpc.Server, error) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()

	lis, err := net.Listen("tcp", ":"+app.ENV.GrpcPort)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcprometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpcprometheus.StreamServerInterceptor),
	)
	grpcHandlers, err := handlers.NewServer(app)
	if err != nil {
		l.Warn("failed to init server").Msg()
		os.Exit(1)
	}
	proto.RegisterLlmServiceServer(server, grpcHandlers)
	reflection.Register(server)

	// init metric for all register
	middleware.GrpcMetric.InitializeMetrics(server)

	go func() {
		l.Infof("🚀 GRPC Server is running on port %s\n", app.ENV.GrpcPort).Msg()
		if err := server.Serve(lis); err != nil {
			l.Warn("GRPC Server failed to start").Msg()
			os.Exit(1)
		}
	}()

	return server, nil
}

func startHTTPGatewayServer(app *bootstrap.Application) (*http.Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()

	server, err := handlers.NewServer(app)
	if err != nil {
		l.Warn("failed to init server").Msg()
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
		l.Warn("Failed to create http gateway").Msg()
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/llm/", http.StripPrefix("/api/v1/llm", grpcMux))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l.Info("404 --- route not found").Msg()
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
		l.Infof("🚀 Server is running on port %v", app.ENV.AppPort).Msg()
		if err := serve.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Warn("server failed to start").Msg()
			os.Exit(1)
		}
	}()

	return serve, nil
}
