package client

import (
	"fmt"
	"testing"
)

func TestUploadFile(t *testing.T) {

	// Create client
	c, err := NewClient(ConnectionConfig{
		Address:   "localhost:9089",
		ChunkSize: 10,
		Compress:  false,
	})

	fmt.Println(c, err)
}
