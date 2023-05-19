package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"io/fs"
	"log"
	"net/http"
)

func GetRoutes(app *app.App) {
	getController := controllers.GetController{
		App: app,
	}

	staticFS, err := fs.Sub(app.Res, "static")
	if err != nil {
		log.Println(err)
		return
	}
	staticHandler := http.FileServer(http.FS(staticFS))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	log.Println("Serving static files from embedded file system /static")

	http.HandleFunc("/", getController.ShowHome)
	http.HandleFunc("/login", getController.ShowLogin)
	http.HandleFunc("/register", getController.ShowRegister)
	http.HandleFunc("/logout", getController.Logout)
}
