package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/bootstrap"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/handlers"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/internal/generated"
	sharedmiddleware "github.com/afrizalsebastian/llm-integration-service/modules/shared-middleware"
	"github.com/gorilla/mux"
)

const apiPrefix = "/api/v1"

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.New()

	app := bootstrap.NewApp()

	mainRouter := mux.NewRouter()
	mainRouter.Use(sharedmiddleware.RecoveryMiddleware())
	mainRouter.Use(sharedmiddleware.CORSMiddleware)

	// Not found handler
	mainRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("404 --- route not found").Msg()
		http.Error(w, "404 - route not found", http.StatusNotFound)
	})

	apiV1Router := mux.NewRouter().PathPrefix(apiPrefix).Subrouter()
	apiServer := handlers.NewServer(app)

	// check app readiness
	apiV1Router.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Service is ready"))
	}).Methods("GET")

	oApiHandlers := generated.HandlerWithOptions(apiServer, generated.GorillaServerOptions{
		BaseRouter: apiV1Router,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request"}`))
		},
	})

	mainRouter.PathPrefix(apiPrefix).Handler(registerCommonMiddleware(
		app,
		oApiHandlers,
	))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", app.ENV.AppPort),
		Handler: mainRouter,
	}

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// start server
	go func() {
		l.Infof("🚀 Server is running on port %v", app.ENV.AppPort).Msg()
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("server failed to start").Msg()
			os.Exit(1)
		}
	}()

	// wait for signal
	<-signalChan
	l.Info("Shutting down server...").Msg()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Error("Server forced to shutdown").Msg()
	}

	l.Info("✅ Server exited gracefully").Msg()
}

func registerCommonMiddleware(_ *bootstrap.Application, handler http.Handler) http.Handler {
	middleware := []sharedmiddleware.MiddlewareFunc{
		sharedmiddleware.RequestTracing(),
	}

	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}
