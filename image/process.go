package image

import (
	"fmt"
	"image"
)

// DIMImage An image.Image extention for pre-processing
type DIMImage struct {
	image.Image
}

//
// Process Tools
//

// ProcessOptions Available options when fetching images
type ProcessOptions struct {
	targetHeight  int64
	targetWidth   int64
	targetQuality float64 // 0.0-1.0
}

// GetTargetSuffix Create filename suffix based on process options
func (i *DIMImage) GetTargetSuffix(opts *ProcessOptions) string {
	return fmt.Sprintf(`%vx%v,%v`, opts.targetWidth, opts.targetHeight, opts.targetQuality)
}

// Resize Resize image based on process options
func (i *DIMImage) Resize(opts *ProcessOptions) (image.Image, error) {

	return i, nil
}
