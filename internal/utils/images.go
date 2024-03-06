package utils

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io"
)

func ResizeImage(reader io.Reader) ([]byte, error) {
	oldImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	var x, y uint
	if oldImage.Bounds().Max.X < oldImage.Bounds().Max.Y {
		y = 480
	} else {
		x = 480
	}

	newImage := resize.Resize(x, y, oldImage, resize.Lanczos3)

	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func ResizeNormalImage(reader io.Reader) ([]byte, error) {
	oldImage, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	var x, y uint
	if oldImage.Bounds().Max.X < oldImage.Bounds().Max.Y {
		y = uint(oldImage.Bounds().Max.Y)
		if oldImage.Bounds().Max.Y > 1920 {
			y = 1920
		}
	} else {
		x = uint(oldImage.Bounds().Max.X)
		if oldImage.Bounds().Max.X > 1920 {
			x = 1920
		}
	}

	newImage := resize.Resize(x, y, oldImage, resize.Lanczos3)

	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, newImage, nil)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
