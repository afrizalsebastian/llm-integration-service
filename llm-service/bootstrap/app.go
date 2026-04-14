package bootstrap

import (
	"context"
	"os"

	appconfig "github.com/afrizalsebastian/go-common-modules/app-config"
	"github.com/afrizalsebastian/go-common-modules/logger"
	"github.com/afrizalsebastian/llm-integration-service/modules/config"
	geminiclient "github.com/afrizalsebastian/llm-integration-service/modules/gemini-client"
)

type Application struct {
	ENV          *config.Config
	GeminiClient geminiclient.IGeminiClient
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}

	l := logger.New()

	wd, err := os.Getwd()
	if err != nil {
		l.Error("Failed to get working directory").Msg()
		return nil
	}

	app.ENV, err = appconfig.Init[config.Config](wd)
	if err != nil {
		l.Error("failed to initialize configuration").Msg()
		os.Exit(1)
	}

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		l.Error("failed to init gemini client").Msg()
		os.Exit(1)
	}
	app.GeminiClient = geminiCient

	return app
}
