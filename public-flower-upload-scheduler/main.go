package main

import (
	"public-flower-upload-scheduler/config"
	"public-flower-upload-scheduler/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	config.ReadConfig()

	e := echo.New()
	e.GET("/", handlers.PingPong)
	e.POST("/upload", handlers.UploadFlowers)
	e.Logger.Fatal(e.Start(":" + config.C.Server.Port))
}
