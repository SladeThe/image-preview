package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"image-preview/internal/app/image-preview/controllers"
)

func init() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetLevel(log.INFO)
}

func main() {
	server := configureWebServer()

	go func() {
		server.Logger.Fatal(server.Start(":9090"))
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Errorf("Unable to stop HTTP server: %s.", err.Error())
	} else {
		server.Logger.Info("HTTP server has exited gracefully.")
	}
}

func configureWebServer() *echo.Echo {
	e := echo.New()

	e.Logger.SetOutput(os.Stdout)
	e.Logger.SetLevel(log.INFO)

	var errorHandler func(err error, c echo.Context)
	if e.HTTPErrorHandler == nil {
		errorHandler = e.DefaultHTTPErrorHandler
	} else {
		errorHandler = e.HTTPErrorHandler
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Warn("Unhandled error:", err)
		errorHandler(err, c)
	}

	e.Use(middleware.Logger())

	controllers.Configure(e)

	return e
}
