package handlers

import (
	"github.com/afrizalsebastian/llm-integration-service/auth-service/application/controllers"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/application/services"
	"github.com/afrizalsebastian/llm-integration-service/auth-service/bootstrap"
)

type ServeController struct {
	GoogleAuthController controllers.IGoogleAuthController
}

func initDI(app *bootstrap.Application) *ServeController {
	return &ServeController{
		GoogleAuthController: googleAuth(app),
	}
}

func googleAuth(app *bootstrap.Application) controllers.IGoogleAuthController {
	googleAuthService := services.NewGoogleAuthService(app.GoogleAuthConfig)
	googleAuthController := controllers.NewGoogleAuthController(googleAuthService)

	return googleAuthController
}
