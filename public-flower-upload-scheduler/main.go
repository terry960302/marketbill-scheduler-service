package main

import (
	"fmt"
	"net/http"
	"public-flower-upload-scheduler/config"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/services"
	"public-flower-upload-scheduler/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	db := datastore.NewPostgresql()
	fmt.Println(db.Config)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + config.C.Server.Port))

	items, err := services.FetchFlowerItems()
	if err != nil {
		panic(err)
	}

	utils.PrettyPrint(items)
}
