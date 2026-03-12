package bootstrap

import (
	"context"
	"log"
	"os"

	appconfig "github.com/afrizalsebastian/llm-integration-service/modules/app-config"
	geminiclient "github.com/afrizalsebastian/llm-integration-service/modules/gemini-client"
)

type Application struct {
	ENV          *appconfig.Config
	GeminiClient geminiclient.IGeminiClient
}

func NewApp() *Application {
	ctx := context.Background()
	app := &Application{}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory")
		return nil
	}

	app.ENV, err = appconfig.Init(wd)
	if err != nil {
		log.Fatal("failed to initialize configuration")
	}

	// Init Gemini Client
	geminiCient, err := geminiclient.NewGeminiAiCLient(ctx, app.ENV.GeminiApiKey, app.ENV.GeminiModel)
	if err != nil {
		log.Fatal("failed to init gemini client")
	}
	app.GeminiClient = geminiCient

	return app
}
