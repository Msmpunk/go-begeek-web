package main

import (
	"GoWeb/app"
	"GoWeb/config"
	"GoWeb/database"
	"GoWeb/models"
	"GoWeb/routes"
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed templates static
var res embed.FS

func main() {
	appLoad := app.App{}

	appLoad.Config = config.LoadConfig()

	appLoad.Res = &res

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Println("Error al crear el archivo")
			log.Panicln(err)
			return
		}
	}

	file, err := os.OpenFile("logs/"+time.Now().Format("2023-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)

	//conectar base de datos

	appLoad.Db = database.ConnectDB(&appLoad)

	if appLoad.Config.Db.AutoMigrate {
		err = models.RunAllMigrations(&appLoad)
		if err != nil {
			log.Println(err)

			return
		}

	}

	appLoad.ScheduledTasks = app.Scheduled{
		EveryReboot: []func(app *app.App){models.SessionCleanUp},
		EveryMinute: []func(app *app.App){models.SessionCleanUp},
	}

	routes.GetRoutes(&appLoad)
	routes.PostRoutes(&appLoad)

	server := &http.Server{Addr: appLoad.Config.Listen.Ip + ":" + appLoad.Config.Listen.Port}

	go func() {
		log.Println("Servidor iniciando")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalln("No se pudo inicializar el server")
		}
	}()

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	stop := make(chan struct{})

	go app.RunScheduledTasks(&appLoad, 100, stop)

	<-interrupt

	log.Println("Se obtuvo una señal de interrupcion, apagando el servidor")

	err = server.Shutdown(context.Background())

	if err != nil {
		log.Fatalln("No se pudo cerrar el servidor ")
	}
}
