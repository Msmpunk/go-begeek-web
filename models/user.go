package models

import (
	"GoWeb/app"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	UserMame  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const userColumnsNoId = "\"Username\", \"Password\", \"CreatedAt\", \"UpdatedAt\""
const userColumns = "\"Id\", " + userColumnsNoId
const userTable = "public.\"User\""

const (
	selectUserById       = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Id\" = $1"
	selectUserByUsername = "SELECT " + userColumns + " FROM " + userTable + " WHERE \"Username\" = $1"
	insertUser           = "INSERT INTO " + userTable + " (" + userColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
)

func GetCurrentUser(app *app.App, r *http.Request) (User, error) {

	cookie, err := r.Cookie("session")

	if err != nil {
		log.Println("No se encontro la seción")
		return User{}, err

	}
	session, err := GetSessionByAuthToken(app, cookie.Value)

	if err != nil {
		log.Println("Error al obtener el token de sesion")
		return User{}, err

	}

	return GetUserById(app, session.UserId)
}

func GetUserById(app *app.App, id int64) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserById, id).Scan(&user.Id, &user.UserMame, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Println("Usuario no encontrado por el id: " + strconv.FormatInt(id, 10))
		return User{}, err

	}
	return user, nil

}

func GetUserByUsername(app *app.App, username string) (User, error) {
	user := User{}

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.UserMame, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error al obtener el usuario:" + username)
		return User{}, err
	}

	return user, nil
}

func CreateUser(app *app.App, username string, password string, createdAt time.Time, updatedAt time.Time) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error encriptando la contraseña")
		return User{}, err
	}

	var lastInsertId int64

	err = app.Db.QueryRow(insertUser, username, string(hash), createdAt, updatedAt).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error creando al usuario")
		return User{}, err
	}

	return GetUserById(app, lastInsertId)
}

func AuthenticateUser(app *app.App, w http.ResponseWriter, username string, password string, remember bool) (Session, error) {
	var user User

	err := app.Db.QueryRow(selectUserByUsername, username).Scan(&user.Id, &user.UserMame, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error de autenticacion, usaurio no encontrado:" + username)
		return Session{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		log.Println("Error de autenticacion,(contraseña no encontrada) para:" + username)
		return Session{}, err
	} else {
		return CreateSession(app, w, user.Id, remember)
	}
}

func LogoutUser(app *app.App, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("Error gobteniendo la cookie")
		return
	}

	err = DeleteSessionByAuthToken(app, w, cookie.Value)
	if err != nil {
		log.Println("Error eliminando la sesion del AuthToken")
		return
	}
}
