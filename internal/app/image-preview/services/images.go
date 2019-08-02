package services

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/nfnt/resize"

	"image-preview/pkg/models/api/requests"
	"image-preview/pkg/models/api/responses"
)

const (
	ImagePreviewWidth  uint = 100
	ImagePreviewHeight uint = 100
)

type Image struct {
}

func (i Image) GeneratePreviews(request *requests.ImagePreview) (response *responses.ImagePreview, err error) {
	var previews []string

	for _, rawImg := range request.Images {
		var imgBytes []byte

		if govalidator.IsURL(rawImg) {
			r, err := http.Get(rawImg)
			if err != nil {
				return nil, err
			}

			if imgBytes, err = ioutil.ReadAll(r.Body); err != nil {
				return nil, err
			}
		} else {
			if imgBytes, err = base64.StdEncoding.DecodeString(rawImg); err != nil {
				return nil, err
			}
		}

		img, err := BytesToImage(imgBytes)
		if err != nil {
			return nil, err
		}

		previewBytes, err := ImageToBytes(resize.Resize(ImagePreviewWidth, ImagePreviewHeight, img, resize.Lanczos3))
		if err != nil {
			return nil, err
		}

		previews = append(previews, base64.StdEncoding.EncodeToString(previewBytes))
	}

	return &responses.ImagePreview{ImageCollection: responses.ImageCollection{Images: previews}}, nil
}

func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func BytesToImage(imgBytes []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	return img, err
}
