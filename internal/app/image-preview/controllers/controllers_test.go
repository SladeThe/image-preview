package controllers

import (
	"bytes"
	"encoding/json"
	"image"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"image-preview/internal/app/image-preview/commons"
	"image-preview/internal/app/image-preview/services"
	"image-preview/internal/app/image-preview/testutils"
	"image-preview/pkg/models/api/requests"
	"image-preview/pkg/models/api/responses"
)

const (
	expectedWidth  = int(services.ImagePreviewWidth)
	expectedHeight = int(services.ImagePreviewHeight)
)

func TestGeneratePreviews_EmptyBody(t *testing.T) {
	e := echo.New()
	Configure(e)

	request := httptest.NewRequest(http.MethodPost, "/api/images/previews", nil)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	if assert.NoError(t, generatePreviews(c)) {
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	}
}

func TestGeneratePreviews_NotAnImage(t *testing.T) {
	e := echo.New()
	Configure(e)

	request := httptest.NewRequest(http.MethodPost, "/api/images/previews", bytes.NewReader([]byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	}))
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	if assert.NoError(t, generatePreviews(c)) {
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	}
}

func TestGeneratePreviews_ImageBlob(t *testing.T) {
	e := echo.New()
	Configure(e)

	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	previewRequest := requests.ImagePreview{ImageCollection: requests.ImageCollection{
		Images: testutils.ImageToBase64SliceOrFail(t, img),
	}}

	jsonRequest, err := json.Marshal(previewRequest)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/images/previews", bytes.NewReader(jsonRequest))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set(echo.HeaderAccept, "image/png")
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.Set(commons.ContextKeyServices, services.Services{})

	if assert.NoError(t, generatePreviews(c)) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		imgBytes, err := ioutil.ReadAll(recorder.Body)
		if err != nil {
			t.Fatal(err)
		}

		img, err := services.BytesToImage(imgBytes)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expectedWidth, img.Bounds().Max.X)
		assert.Equal(t, expectedHeight, img.Bounds().Max.Y)
	}
}

func TestGeneratePreviews_ImageJSON(t *testing.T) {
	e := echo.New()
	Configure(e)

	img := image.NewRGBA(image.Rect(0, 0, 500, 500))
	previewRequest := requests.ImagePreview{ImageCollection: requests.ImageCollection{
		Images: testutils.ImageToBase64SliceOrFail(t, img),
	}}

	jsonRequest, err := json.Marshal(previewRequest)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/images/previews", bytes.NewReader(jsonRequest))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)
	c.Set(commons.ContextKeyServices, services.Services{})

	if assert.NoError(t, generatePreviews(c)) {
		assert.Equal(t, http.StatusOK, recorder.Code)

		jsonResponse, err := ioutil.ReadAll(recorder.Body)
		if err != nil {
			t.Fatal(err)
		}

		var response responses.ImagePreview
		err = json.Unmarshal(jsonResponse, &response)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 1, len(response.Images))

		img := testutils.Base64ToImageOrFail(t, response.Images[0])

		assert.Equal(t, expectedWidth, img.Bounds().Max.X)
		assert.Equal(t, expectedHeight, img.Bounds().Max.Y)
	}
}
