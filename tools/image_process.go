package tools

import "fmt"

// DIMImage is an abstracted image type for pre-processing
type DIMImage interface {
	GetTargetSuffix() string
}

// ImageProcessOptions contains available options when fetching images
type ImageProcessOptions struct {
	targetHeight  int64
	targetWidth   int64
	targetQuality float64 // 0.0-1.0
}

// GetTargetSuffix create filename suffix based on process options
func (o *ImageProcessOptions) GetTargetSuffix() string {
	return fmt.Sprintf(`%bx%b,%b`, o.targetWidth, o.targetWidth, o.targetQuality)
}
