package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HttpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
	c.String(http.StatusOK, err.Error())
}
