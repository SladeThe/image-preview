package middleware

import (
	"github.com/labstack/echo/v4"

	"image-preview/internal/app/image-preview/commons"
	"image-preview/internal/app/image-preview/services"
)

func SetServices(services services.Services) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(commons.ContextKeyServices, services)
			return next(c)
		}
	}
}
