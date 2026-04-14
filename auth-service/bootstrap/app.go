package bootstrap

import (
	"os"

	appconfig "github.com/afrizalsebastian/go-common-modules/app-config"
	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/modules/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Application struct {
	ENV              *config.Config
	GoogleAuthConfig *oauth2.Config
}

func NewApp() *Application {
	app := &Application{}

	l := logger.New()

	wd, err := os.Getwd()
	if err != nil {
		l.Error("Failed to get working directory").Msg()
		os.Exit(1)
	}

	app.ENV, err = appconfig.Init[config.Config](wd)
	if err != nil {
		l.Error("failed to initialize configuration").Msg()
		os.Exit(1)
	}

	googleOauthConfig := &oauth2.Config{
		ClientID:     app.ENV.GoogleAuthClientID,
		ClientSecret: app.ENV.GoogleAuthClientSecret,
		RedirectURL:  app.ENV.GoogleCallbackUrl,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	app.GoogleAuthConfig = googleOauthConfig

	return app
}
