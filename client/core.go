package client

import (
	"context"
	"dim-fs/protocol"
	"dim-fs/utils"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// ConnectionInstance conatins grpc connection instance
type ConnectionInstance struct {
	conn      *grpc.ClientConn
	client    protocol.CoreServiceClient
	chunkSize int
}

// ConnectionConfig contains basic config
type ConnectionConfig struct {
	Address   string
	ChunkSize int
	Compress  bool
}

// Stat benchmark
type Stat struct {
	StartedAt  time.Time
	FinishedAt time.Time
}

// NewClient create new grpc client
func NewClient(cfg ConnectionConfig) (c ConnectionInstance, err error) {
	var (
		grpcOpts = []grpc.DialOption{}
	)

	if cfg.Address == "" {
		err = errors.Errorf("Address must be specified")
		return
	}

	if cfg.Compress {
		grpcOpts = append(grpcOpts,
			grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	grpcOpts = append(grpcOpts, grpc.WithInsecure())

	switch {
	case cfg.ChunkSize == 0:
		err = errors.Errorf("ChunkSize must be specified")
		return
	case cfg.ChunkSize > (1 << 22):
		err = errors.Errorf("ChunkSize must be < than 4MB")
		return
	default:
		c.chunkSize = cfg.ChunkSize
	}

	c.conn, err = grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		err = errors.Wrapf(err,
			"Failed to start grpc connection with address %s",
			cfg.Address)
		return
	}

	c.client = protocol.NewCoreServiceClient(c.conn)

	return
}

// UploadFile upload file handler
func (c *ConnectionInstance) UploadFile(ctx context.Context, f string) (stats Stat, err error) {
	var (
		writing  = true
		buf      []byte
		n        int
		file     *os.File
		response *protocol.UploadFileResponse
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err,
			"Failed to open file %s",
			f)
		return
	}
	defer file.Close()

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		err = errors.Wrapf(err,
			"Failed to create upload stream for file %s",
			f)
		return
	}
	defer stream.CloseSend()

	stats.StartedAt = time.Now()
	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {

				writing = false
				err = nil
				continue
			}

			err = errors.Wrapf(err,
				"Error while copying from file to buf")
			return
		}

		err = stream.Send(&protocol.Chunk{
			Content:  buf[:n],
			FileName: file.Name() + uuid.New().String(),
		})
		if err != nil {
			err = errors.Wrapf(err,
				"Failed to send chunk via stream")
			return
		}
	}

	stats.FinishedAt = time.Now()

	response, err = stream.CloseAndRecv()
	fmt.Println(buf)

	if err != nil {
		err = errors.Wrapf(err,
			"Failed to receive upstream status response")
		return
	}

	if response.UploadStatusCode != protocol.UploadStatusCode_Ok {
		err = errors.Errorf(
			"Upload failed - msg: %s",
			response.Message)
		return
	}

	return
}

// DownloadFile download file handler
func (c *ConnectionInstance) DownloadFile(ctx context.Context, fileName string) (err error) {
	var f *os.File
	var chunks *protocol.ResChunk
	stream, err := c.client.DownloadFile(ctx, &protocol.DownloadFileParams{FileName: fileName})

	for {

		// Get chunks from stream
		chunks, err = stream.Recv()

		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"Failed unexpectadely while reading chunks from stream")
			return
		}

		// Create file
		f, err = os.Create(viper.GetString("file.upload") + fileName)

		if err != nil {
			utils.LogError("Unable to create file")

			return
		}
		defer f.Close()

		// Write into file
		err = utils.WriteToFile(f, chunks.Content)
		if err != nil {
			utils.LogError("Unable to write chunk of filename :" + err.Error())

			return
		}
	}

	if err != nil {
		err = errors.Wrapf(err,
			"Failed to send status code")
		return
	}

	return
}

// Close close grpc connection
func (c *ConnectionInstance) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
