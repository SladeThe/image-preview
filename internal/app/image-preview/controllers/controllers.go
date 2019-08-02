package controllers

import (
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/labstack/echo/v4"

	"image-preview/internal/app/image-preview/commons"
	"image-preview/internal/app/image-preview/middleware"
	"image-preview/internal/app/image-preview/services"
	"image-preview/internal/app/image-preview/validation"
	"image-preview/pkg/models/api/requests"
)

func Configure(e *echo.Echo) {
	e.Use(middleware.SetServices(services.Services{
		Image: services.Image{},
	}))

	apiGroup := e.Group("/api")

	apiGroup.POST("/images/previews", generatePreviews)
}

func generatePreviews(c echo.Context) error {
	var request requests.ImagePreview
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if form, err := c.MultipartForm(); err == nil {
		for _, imgHeader := range form.File["image"] {
			img, err := func(imgHeader *multipart.FileHeader) ([]byte, error) {
				img, err := imgHeader.Open()
				if err != nil {
					return nil, err
				}
				defer img.Close()

				return ioutil.ReadAll(img)
			}(imgHeader)

			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			request.Images = append(request.Images, base64.StdEncoding.EncodeToString(img))
		}
	}

	if err := validation.Validate(request); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	response, err := c.Get(commons.ContextKeyServices).(services.Services).Image.GeneratePreviews(&request)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if c.Request().Header.Get(echo.HeaderAccept) == "image/png" && len(response.Images) == 1 {
		imgBytes, err := base64.StdEncoding.DecodeString(response.Images[0])
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.Blob(http.StatusOK, "image/png", imgBytes)
	}

	return c.JSON(http.StatusOK, response)
}
