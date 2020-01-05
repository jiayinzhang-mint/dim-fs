package image

import (
	"image"
	"image/jpeg"
	"io"
	"path/filepath"
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

// Type get image type from path
func Type(p string) string {
	extname := filepath.Ext(p)
	switch extname {
	case ".jpg", ".jpeg", ".jpe", ".jfif":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".ico":
		return "image/x-icon"
	case ".gif":
		return "image/gif"
	case ".tif", ".tiff":
		return "image/jiff"
	default:
		return "text/plain"
	}
}
