package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"GoWeb/templating"
	"log"
	"net/http"
)

type GetCotroller struct {
	App *app.App
}

func (getCotroller *GetCotroller) ShowHome(w http.ResponseWriter, _ *http.Request) {
	type dataStruct struct {
		Test string
	}
	data := dataStruct{
		Test: "Hola mundo",
	}

	templating.RenderTemplate(getCotroller.App, w, "templates/pages/home.html", data)
}

func (getCotroller *GetCotroller) ShowRegister(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		tokencsrf string
	}
	tokencsrf, err := security.GenerateToken(w, r)

	if err != nil {
		log.Println("Error generarndo el token")
		return
	}

	data := dataStruct{
		tokencsrf: tokencsrf,
	}

	templating.RenderTemplate(getCotroller.App, w, "templates/pages/register.html", data)
}

func (getCotroller *GetCotroller) ShowLogin(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		tokencsrf string
	}
	tokencsrf, err := security.GenerateToken(w, r)

	if err != nil {
		log.Println("Error generarndo el token")
		return
	}

	data := dataStruct{
		tokencsrf: tokencsrf,
	}

	templating.RenderTemplate(getCotroller.App, w, "templates/pages/login.html", data)
}

func (getCotroller *GetCotroller) Logout(w http.ResponseWriter, r *http.Request) {

	models.LogoutUser(getCotroller.App, w, r)
	http.Redirect(w, r, "/", http.StatusNotFound)
}
