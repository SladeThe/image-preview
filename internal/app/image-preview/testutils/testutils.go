package testutils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"
)

func ImageToBase64SliceOrFail(t *testing.T, img image.Image) []string {
	buf := new(bytes.Buffer)

	if err := png.Encode(buf, img); err != nil {
		t.Fatal(err)
	}

	return append(make([]string, 0, 1), base64.StdEncoding.EncodeToString(buf.Bytes()))
}

func Base64ToImageOrFail(t *testing.T, imgBase64 string) image.Image {
	imgBytes, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		t.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		t.Fatal(err)
	}

	return img
}
