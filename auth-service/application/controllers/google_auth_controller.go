package controllers

import (
	"context"
	"net/http"

	"log"

	"github.com/afrizalsebastian/llm-integration-service/auth-service/application/services"
)

type IGoogleAuthController interface {
	GoogleOauthRedirect(context.Context, http.ResponseWriter, *http.Request) string
	GoogleAuthCallback(context.Context, http.ResponseWriter, *http.Request) string
}

type googleAuthContoller struct {
	GoogleAuthService services.IGoogleAuthService
}

func NewGoogleAuthController(googleAuthService services.IGoogleAuthService) IGoogleAuthController {
	return &googleAuthContoller{
		GoogleAuthService: googleAuthService,
	}
}

func (g *googleAuthContoller) GoogleOauthRedirect(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	return g.GoogleAuthService.GoogleOauthRedirect(ctx, w)
}

func (g *googleAuthContoller) GoogleAuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request) string {
	cookieState, err := r.Cookie("oauth_state")
	if err != nil || cookieState.Value != r.URL.Query().Get("state") {
		log.Println("invalid oauth_state")
		return "/not-found"
	}

	http.SetCookie(w, &http.Cookie{Name: "oauth_state", MaxAge: -1})

	code := r.URL.Query().Get("code")

	return g.GoogleAuthService.GoogleAuthCallback(ctx, cookieState, code)
}
