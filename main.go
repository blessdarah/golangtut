package main

import (
	"blessdarah/tuts/user"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type api struct {
	store user.Store
}

func (a *api) listUsers(c *echo.Context) error {
	return c.JSON(http.StatusOK, a.store.GetAll())
}

func (a *api) createUser(c *echo.Context) error {
	var user user.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	err := user.Validate()
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]any{"errors": err},
		)
	}

	a.store.Add(user)

	return c.JSON(http.StatusCreated, user)
}

func main() {
	e := echo.New()

	api := api{
		store: user.NewUserStore(),
	}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/users", api.listUsers)
	e.POST("/users", api.createUser)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
