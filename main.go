package main

import (
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	TestResponse struct {
		Return string `json:"return" xml:"return"`
	}
)

func NewRouter() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/", Version)

	return e
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}

// Default Handler
func Version(c echo.Context) error {
	// No malicious links are found
	return c.JSON(http.StatusOK, &TestResponse{
		Return: "Hello World",
	})
}
