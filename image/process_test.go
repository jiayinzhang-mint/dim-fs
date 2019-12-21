package image

import (
	"fmt"
	"testing"
)

func TestGetTargetSuffix(*testing.T) {
	var i DIMImage
	opts := &ProcessOptions{
		targetHeight:  1,
		targetWidth:   2,
		targetQuality: 3,
	}
	fmt.Println(i.GetTargetSuffix(opts))
}
