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

	stat, uploadFileErr := c.UploadFile(context.TODO(), "test.txt", "biu")
	fmt.Println(stat, uploadFileErr)
}

func TestDownloadFile(t *testing.T) {
	// Create client
	c, createClientErr := NewClient(ConnectionConfig{
		Address:   "localhost:9089",
		ChunkSize: 10,
		Compress:  false,
	})
	fmt.Println(c, createClientErr)

	err := c.DownloadFile(context.TODO(), "63f826b0-9f2b-4cc0-af2c-c73f5cf1a8e9.test.txt", "")
	fmt.Println(err)
}
