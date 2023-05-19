package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"GoWeb/templating"
	"fmt"
	"net/http"
)

type GetController struct {
	App *app.App
}

func (getController *GetController) ShowHome(w http.ResponseWriter, _ *http.Request) {
	type dataStruct struct {
		Test string
	}

	data := dataStruct{
		Test: "Hola Mundo",
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/home.html", data)
}

func (getController *GetController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	CsrfToken, err := security.GenerateToken(w, r)
	if err != nil {
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/register.html", data)
}

func (getController *GetController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	type dataStruct struct {
		CsrfToken string
	}

	CsrfToken, err := security.GenerateToken(w, r)
	fmt.Println(CsrfToken)
	if err != nil {
		return
	}

	data := dataStruct{
		CsrfToken: CsrfToken,
	}

	templating.RenderTemplate(getController.App, w, "templates/pages/login.html", data)
}

func (getController *GetController) Logout(w http.ResponseWriter, r *http.Request) {
	models.LogoutUser(getController.App, w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}
