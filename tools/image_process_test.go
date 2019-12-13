package tools

import (
	"fmt"
	"testing"
)

func TestGetTargetSuffix(*testing.T) {
	var i DIMImage
	opts := &ImageProcessOptions{
		targetHeight:  1,
		targetWidth:   2,
		targetQuality: 3,
	}
	fmt.Println(i.GetTargetSuffix(opts))
}
