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

	stat, uploadFileErr := c.UploadFile(context.TODO(), "test.txt")
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

	err := c.DownloadFile(context.TODO(), "c8891a39-314c-4af9-a814-7e461fc61972.test.txt")
	fmt.Println(err)
}
