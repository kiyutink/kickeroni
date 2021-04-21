package server

import (
	"github.com/labstack/echo/v4"
)

func Start() {
	e := echo.New()
	e.POST("/events", Events)
	e.POST("/play", Play)
	e.POST("/ranks", Ranks)
	e.POST("/interactions", Interactions)
	e.HTTPErrorHandler = HttpErrorHandler
	e.Logger.Fatal(e.Start(":3000"))
}
