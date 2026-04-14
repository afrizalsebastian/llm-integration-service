package bootstrap

import (
	"log"
	"os"

	appconfig "github.com/afrizalsebastian/go-common-modules/app-config"
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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory")
	}

	app.ENV, err = appconfig.Init[config.Config](wd)
	if err != nil {
		log.Fatal("failed to initialize configuration")
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
