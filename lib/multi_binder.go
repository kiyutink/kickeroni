package lib

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
)

type MultiBinder struct{}

func (cb *MultiBinder) Bind(i interface{}, c echo.Context) (err error) {
	body, err := io.ReadAll(c.Request().Body)
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	db := new(echo.DefaultBinder)
	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
		c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
		return
	}

	return err
}
