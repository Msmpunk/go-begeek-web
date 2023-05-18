package controllers

import (
	"GoWeb/app"
	"GoWeb/models"
	"GoWeb/security"
	"log"
	"net/http"
	"time"
)

type PostController struct {
	App *app.App
}

func (postController *PostController) Login(w http.ResponseWriter, r *http.Request) {

	_, err := security.VerifyToken(r)

	if err != nil {
		log.Println("Error verficando el csrf token ")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	rememberme := r.FormValue("remember") == "on"

	if username == "" || password == "" {
		log.Println("El usuario o la contraseña estan  vacios")
		http.Redirect(w, r, "/login", http.StatusNotFound)
	}

	_, err = models.AuthenticateUser(postController.App, w, username, password, rememberme)

	if err != nil {
		log.Println("usuario no valido")
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusNotFound)
	}

	http.Redirect(w, r, "/", http.StatusNotFound)
}

func (postController *PostController) Register(w http.ResponseWriter, r *http.Request) {

	_, err := security.VerifyToken(r)

	if err != nil {
		log.Println("Error verficando el csrf token ")
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	createdAt := time.Now()
	updatedAt := time.Now()

	if username == "" || password == "" {
		log.Println("El usuario o la contraseña estan  vacios")
		http.Redirect(w, r, "/register", http.StatusNotFound)
	}

	_, err = models.CreateUser(postController.App, username, password, createdAt, updatedAt)

	if err != nil {
		log.Println("Error al crear usuario ")
		log.Println(err)
	}

	http.Redirect(w, r, "/login", http.StatusNotFound)
}
