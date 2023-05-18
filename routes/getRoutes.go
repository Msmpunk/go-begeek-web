package routes

import (
	"GoWeb/app"
	"GoWeb/controllers"
	"io/fs"
	"log"
	"net/http"
)

func GetRoutes(app *app.App) {

	getCotroller := controllers.GetCotroller{
		App: app,
	}

	staticFS, err := fs.Sub(app.Res, "static")

	if err != nil {
		log.Println(err)
		return
	}

	staticHandler := http.FileServer(http.FS(staticFS))

	http.Handle("/static", http.StripPrefix("/static", staticHandler))

	log.Println("Archivos cargados correctamente")

	//pages
	http.HandleFunc("/", getCotroller.ShowHome)
	http.HandleFunc("/login", getCotroller.ShowLogin)
	http.HandleFunc("/register", getCotroller.ShowRegister)
	http.HandleFunc("/logout", getCotroller.Logout)
}
