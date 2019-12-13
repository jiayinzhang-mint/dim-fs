package tools

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
)

//
// Decode Tools
//

// DIMImage An image.Image extention for pre-processing
type DIMImage struct {
	image.Image
}

// DecodeToJpeg Convert iostream to jpeg image
func DecodeToJpeg(r io.Reader) (DIMImage, error) {
	i, err := jpeg.Decode(r)
	return DIMImage{i}, err
}

//
// Process Tools
//

// ImageProcessOptions Available options when fetching images
type ImageProcessOptions struct {
	targetHeight  int64
	targetWidth   int64
	targetQuality float64 // 0.0-1.0
}

// GetTargetSuffix Create filename suffix based on process options
func (i *DIMImage) GetTargetSuffix(opts *ImageProcessOptions) string {
	return fmt.Sprintf(`%vx%v,%v`, opts.targetWidth, opts.targetHeight, opts.targetQuality)
}

// Resize Resize image based on process options
func (i *DIMImage) Resize(opts *ImageProcessOptions) (image.Image, error) {
	// TODO
	return i, nil
}
