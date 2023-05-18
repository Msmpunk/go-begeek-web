package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"net/http"
)

func PostRoutes(app *app.App) {
	postController := controllers.PostController{
		App: app,
	}

	http.HandleFunc("/register-handle", postController.Register)
	http.HandleFunc("/login-handle", postController.Login)
}
