package database

import (
	"GoWeb/app"
	"database/sql"
	"fmt"
	"log"
)

// funcion que hace la coneccion a la base de datos

func ConnectDB(app *app.App) *sql.DB {

	postgresconfig := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", app.Config.Db.Ip, app.Config.Db.Port, app.Config.Db.User, app.Config.Db.Password, app.Config.Db.Name)

	db, err := sql.Open("postgres", postgresconfig)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	log.Println("Conexion exitosa a nuestra base de datos")

	return db

}
