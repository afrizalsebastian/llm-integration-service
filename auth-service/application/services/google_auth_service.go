package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/domain/dto"
	"golang.org/x/oauth2"
)

type IGoogleAuthService interface {
	GoogleOauthRedirect(context.Context, http.ResponseWriter) string
	GoogleAuthCallback(context.Context, *http.Cookie, string) string
}

type googleAuthService struct {
	GoogleOauthConfig *oauth2.Config
}

func NewGoogleAuthService(googleOauthConfig *oauth2.Config) IGoogleAuthService {
	return &googleAuthService{
		GoogleOauthConfig: googleOauthConfig,
	}
}

func (g *googleAuthService) GoogleOauthRedirect(ctx context.Context, w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(
		w, &http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		},
	)

	url := g.GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return url
}

func (g *googleAuthService) GoogleAuthCallback(ctx context.Context, cookieState *http.Cookie, code string) string {
	l := logger.New().WithContext(ctx)

	token, err := g.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		l.Error("error when exchange code to google").Msg()
		return "/not-found"
	}

	client := g.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		l.Error("error when get profile from google").Msg()
		return "/not-found"
	}

	defer resp.Body.Close()

	var googleUser dto.GoogleUser
	if err = json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		l.Error("error when unmarshal user profile data").Msg()
		return "/not-found"
	}

	return "/api/v1/readiness"
}
