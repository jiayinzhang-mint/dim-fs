package client

import (
	"context"
	"dim-fs/protocol"
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
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
		err = errors.Errorf("address must be specified")
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
			"failed to start grpc connection with address %s",
			cfg.Address)
		return
	}

	c.client = protocol.NewCoreServiceClient(c.conn)

	return
}

// UploadFile upload file handler
func (c *ConnectionInstance) UploadFile(ctx context.Context, f string) (stats Stat, err error) {
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *protocol.UploadFileResponse
	)

	file, err = os.Open(f)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open file %s",
			f)
		return
	}
	defer file.Close()

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create upload stream for file %s",
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
				"errored while copying from file to buf")
			return
		}

		err = stream.Send(&protocol.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			err = errors.Wrapf(err,
				"failed to send chunk via stream")
			return
		}
	}

	stats.FinishedAt = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to receive upstream status response")
		return
	}

	if status.UploadStatusCode != protocol.UploadStatusCode_Ok {
		err = errors.Errorf(
			"upload failed - msg: %s",
			status.Message)
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
