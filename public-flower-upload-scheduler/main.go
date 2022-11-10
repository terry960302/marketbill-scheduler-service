package main

import (
	"fmt"
	"public-flower-upload-scheduler/config"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	db := datastore.NewPostgresql()
	fmt.Println(db.Config)

	e := echo.New()
	e.GET("/", handlers.PingPong)
	e.POST("/upload", handlers.UploadFlowers)
	e.Logger.Fatal(e.Start(":" + config.C.Server.Port))
}
