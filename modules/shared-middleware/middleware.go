package sharedmiddleware

import (
	"context"
	"net/http"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/google/uuid"
)

type MiddlewareFunc func(http.Handler) http.Handler

func RecoveryMiddleware() func(next http.Handler) http.Handler {
	l := logger.New()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					l.Info("Recoverd from panic").Msg()
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func RequestTracing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			msgId := uuid.New().String()
			ctx := context.WithValue(r.Context(), logger.ContextKeyMessageId, msgId)

			w.Header().Set("X-Message-ID", msgId)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
