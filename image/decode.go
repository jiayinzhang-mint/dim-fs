package image

import (
	"image"
	"image/jpeg"
	"io"
)

//
// Decode Tools
//

// DecodeToJpeg Convert iostream to jpeg image
func DecodeToJpeg(r io.Reader) (DIMImage, error) {
	i, err := jpeg.Decode(r)
	return DIMImage{i}, err
}

// DecodeInfo Get image height, width and color model
func DecodeInfo(r io.Reader) (image.Config, error) {
	cfg, _, err := image.DecodeConfig(r)
	return cfg, err
}
