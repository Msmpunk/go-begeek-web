package security

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"math"
	"net/http"
)

// Aqui se generara los tokens de usaurio

func GenerateToken(w http.ResponseWriter, _ *http.Request) (string, error) {

	buff := make([]byte, int(math.Ceil(64)/2))

	_, err := rand.Read(buff)

	if err != nil {

		log.Println("Error creando el buffer")
		log.Panicln(err)

		return "", err

	}

	str := hex.EncodeToString(buff)
	token := str[:64]

	cookie := &http.Cookie{
		Name:     "token_beegek",
		Value:    token,
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	return token, nil
}

func VerifyToken(r *http.Request) (bool, error) {

	cookie, err := r.Cookie("token_beegek")

	if err != nil {

		log.Println("Error tomando el token")
		log.Panicln(err)

		return false, err

	}

	token := r.FormValue("token_beegek")

	if cookie.Value == token {
		return true, nil
	}

	return false, nil

}
