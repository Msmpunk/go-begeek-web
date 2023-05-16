package main

import (
	"GoWeb/app"
	"GoWeb/config"
	"log"
	"os"
)

func main() {
	appLoad := app.App{}

	appLoad.Config = config.LoadConfig()

	// Cargar templates

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Println("Error al crear el archivo")
			log.Panicln(err)

		}
	}

	// file, err := os.OpenFile("logs/"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// log.SetOutput(file)

	// //conectar base de datos

	appLoad.ScheduledTasks = app.Scheduled{
		EveryReboot: []func(app *app.App){},
	}

}
