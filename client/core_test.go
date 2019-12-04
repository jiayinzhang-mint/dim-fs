package client

import (
	"context"
	"fmt"
	"testing"
)

func TestUploadFile(t *testing.T) {

	// Create client
	c, createClientErr := NewClient(ConnectionConfig{
		Address:   "localhost:9089",
		ChunkSize: 10,
		Compress:  false,
	})
	fmt.Println(c, createClientErr)

	stat, uploadFileErr := c.UploadFile(context.TODO(), "./test.txt")
	fmt.Println(stat, uploadFileErr)
}
