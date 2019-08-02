package services

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"testing"

	"gopkg.in/go-playground/assert.v1"

	"image-preview/internal/app/image-preview/testutils"
	"image-preview/pkg/models/api/requests"
)

const (
	expectedWidth  = int(ImagePreviewWidth)
	expectedHeight = int(ImagePreviewHeight)
)

var expectedBounds = image.Rectangle{
	Min: image.Point{X: 0, Y: 0},
	Max: image.Point{X: expectedWidth, Y: expectedHeight},
}

func TestImage_GeneratePreviews_FixedColor(t *testing.T) {
	rand.Seed(4634564536)

	for i := 0; i < 10; i++ {
		width, height := 10+rand.Intn(190), 10+rand.Intn(190)
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		components := make([]byte, 3)
		rand.Read(components)

		fixedColor := color.RGBA{R: components[0], G: components[1], B: components[2], A: math.MaxUint8}

		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				img.Set(x, y, fixedColor)
			}
		}

		images := testutils.ImageToBase64SliceOrFail(t, img)
		request := requests.ImagePreview{ImageCollection: requests.ImageCollection{Images: images}}

		response, err := Image{}.GeneratePreviews(&request)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(response.Images), 1)

		preview := testutils.Base64ToImageOrFail(t, response.Images[0])

		assert.Equal(t, preview.Bounds(), expectedBounds)

		for x := 0; x < expectedWidth; x++ {
			for y := 0; y < expectedHeight; y++ {
				assert.Equal(t, preview.At(x, y), fixedColor)
			}
		}
	}
}

func TestImage_GeneratePreviews_GradientColor(t *testing.T) {
	width, height := 1000, 1000
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA64{
				R: uint16((float64(x) / float64(width)) * math.MaxUint16),
				G: uint16((float64(y) / float64(height)) * math.MaxUint16),
				B: uint16((float64(x+y) / float64(width+height)) * math.MaxUint16),
				A: math.MaxUint16,
			})
		}
	}

	images := testutils.ImageToBase64SliceOrFail(t, img)
	request := requests.ImagePreview{ImageCollection: requests.ImageCollection{Images: images}}

	response, err := Image{}.GeneratePreviews(&request)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(response.Images), 1)

	preview := testutils.Base64ToImageOrFail(t, response.Images[0])

	assert.Equal(t, preview.Bounds(), expectedBounds)

	for x := 0; x < expectedWidth; x++ {
		for y := 0; y < expectedHeight; y++ {
			r, g, b, a := preview.At(x, y).RGBA()

			expectedColor := color.RGBA64{
				R: uint16((float64(x) / float64(expectedWidth)) * math.MaxUint16),
				G: uint16((float64(y) / float64(expectedHeight)) * math.MaxUint16),
				B: uint16((float64(x+y) / float64(expectedWidth+expectedHeight)) * math.MaxUint16),
				A: math.MaxUint16,
			}

			delta := float64(1000.0)

			if math.Abs(float64(int64(r)-int64(expectedColor.R))) > delta {
				t.Fatalf("Expected R %d, but got %d.", expectedColor.R, r)
			}

			if math.Abs(float64(int64(g)-int64(expectedColor.G))) > delta {
				t.Fatalf("Expected G %d, but got %d.", expectedColor.G, g)
			}

			if math.Abs(float64(int64(b)-int64(expectedColor.B))) > delta {
				t.Fatalf("Expected B %d, but got %d.", expectedColor.B, b)
			}

			if math.Abs(float64(int64(a)-int64(expectedColor.A))) > delta {
				t.Fatalf("Expected A %d, but got %d.", expectedColor.A, a)
			}
		}
	}
}
