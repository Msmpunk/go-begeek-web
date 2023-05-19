package models

import (
	"GoWeb/app"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

type Session struct {
	Id         int64
	UserId     int64
	AuthToken  string
	RememberMe bool
	CreatedAt  time.Time
}

const sessionColumnsNoId = "\"UserId\", \"AuthToken\",\"RememberMe\", \"CreatedAt\""
const sessionColumns = "\"Id\", " + sessionColumnsNoId
const sessionTable = "public.\"Session\""

const (
	selectSessionByAuthToken      = "SELECT " + sessionColumns + " FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
	selectAuthTokenIfExists       = "SELECT EXISTS(SELECT 1 FROM " + sessionTable + " WHERE \"AuthToken\" = $1)"
	insertSession                 = "INSERT INTO " + sessionTable + " (" + sessionColumnsNoId + ") VALUES ($1, $2, $3, $4) RETURNING \"Id\""
	deleteSessionByAuthToken      = "DELETE FROM " + sessionTable + " WHERE \"AuthToken\" = $1"
	deleteSessionsOlderThan30Days = "DELETE FROM " + sessionTable + " WHERE \"CreatedAt\" < NOW() - INTERVAL '30 days'"
	deleteSessionsOlderThan6Hours = "DELETE FROM " + sessionTable + " WHERE \"CreatedAt\" < NOW() - INTERVAL '6 hours' AND \"RememberMe\" = false"
)

func CreateSession(app *app.App, w http.ResponseWriter, userId int64, remember bool) (Session, error) {
	session := Session{}
	session.UserId = userId
	session.AuthToken = generateAuthToken(app)
	session.RememberMe = remember
	session.CreatedAt = time.Now()

	var existingAuthToken bool
	err := app.Db.QueryRow(selectAuthTokenIfExists, session.AuthToken).Scan(&existingAuthToken)
	if err != nil {
		log.Println("Error validando los tokens existentes")
		log.Println(err)
		return Session{}, err
	}

	if existingAuthToken == true {
		log.Println("Hay un token duplicado en la tabla sessiones, generando un nuevo token")
		return CreateSession(app, w, userId, remember)
	}

	err = app.Db.QueryRow(insertSession, session.UserId, session.AuthToken, session.RememberMe, session.CreatedAt).Scan(&session.Id)
	if err != nil {
		log.Println("Error al insertar en la base de datos")
		return Session{}, err
	}

	createSessionCookie(app, w, session)
	return session, nil
}

func GetSessionByAuthToken(app *app.App, authToken string) (Session, error) {
	session := Session{}

	err := app.Db.QueryRow(selectSessionByAuthToken, authToken).Scan(&session.Id, &session.UserId, &session.AuthToken, &session.RememberMe, &session.CreatedAt)

	if err != nil {

		log.Println("Error obteniedo la sesion por auth token")
		log.Println(err)
		return Session{}, err

	}
	return session, nil

}

func generateAuthToken(app *app.App) string {
	b := make([]byte, 64)
	_, err := rand.Read(b)

	if err != nil {
		log.Println("Error creando el buffer")
	}

	return hex.EncodeToString(b)
}

func createSessionCookie(app *app.App, w http.ResponseWriter, session Session) {
	cookie := &http.Cookie{}

	if session.RememberMe {
		cookie = &http.Cookie{
			Name:     "session",
			Value:    session.AuthToken,
			Path:     "/",
			MaxAge:   2592000 * 1000, // 30minutos en la session
			HttpOnly: true,
			Secure:   true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     "session",
			Value:    session.AuthToken,
			Path:     "/",
			MaxAge:   21600 * 1000, // 6minutos de tiempo en la session
			HttpOnly: true,
			Secure:   true,
		}
	}

	http.SetCookie(w, cookie)
}

func deleteSessionCookie(app *app.App, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func DeleteSessionByAuthToken(app *app.App, w http.ResponseWriter, authToken string) error {
	_, err := app.Db.Exec(deleteSessionByAuthToken, authToken)

	if err != nil {
		log.Println("Error al borrar de la base de datos")
		return err
	}

	deleteSessionCookie(app, w)
	return nil
}

func SessionCleanUp(app *app.App) {
	_, err := app.Db.Exec(deleteSessionsOlderThan30Days)

	if err != nil {
		log.Println("Error borrando las sesiones de mas de 30 dias")
		log.Println(err)
	}

	_, err = app.Db.Exec(deleteSessionsOlderThan6Hours)

	if err != nil {
		log.Println("Error borrando las sesiones de mas de 6  horas")
		log.Println(err)
	}

	log.Println("Todas las seciones han sido removidas")

}
