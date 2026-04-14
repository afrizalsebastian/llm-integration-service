package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/afrizalsebastian/llm-integration-service/auth-service/application/controllers"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/bootstrap"
)

type Server struct {
	GoogleAuthController controllers.IGoogleAuthController
}

func NewServer(app *bootstrap.Application) *Server {
	di := initDI(app)
	return &Server{
		GoogleAuthController: di.GoogleAuthController,
	}
}

func (s *Server) GetAuthGoogle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.GoogleAuthController.GoogleOauthRedirect(ctx, w, r)
	http.Redirect(w, r, resp, http.StatusTemporaryRedirect)
}

func (s *Server) GetAuthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp := s.GoogleAuthController.GoogleAuthCallback(ctx, w, r)
	http.Redirect(w, r, resp, http.StatusSeeOther)
}
