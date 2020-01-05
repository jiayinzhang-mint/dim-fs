package rpc

import (
	"fmt"
	"github.com/insdim/dim-fs/protocol"
	"github.com/insdim/dim-fs/utils"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CoreService contains core file system service
type CoreService struct{}

// UploadFile handle upload file call
func (c *CoreService) UploadFile(stream protocol.CoreService_UploadFileServer) (err error) {
	firstChunk := true
	var f *os.File
	var chunks *protocol.Chunk

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

		if firstChunk { // First chunk contains file name
			// Check file name
			var fileName string
			if chunks.FileName != "" {
				fileName = chunks.FileName
			} else {
				fileName = uuid.New().String()
			}

			// Create path if not exist
			fileDir := filepath.Dir(fileName)
			if _, err := os.Stat(viper.GetString("file.upload") + fileDir); os.IsNotExist(err) {
				os.MkdirAll(viper.GetString("file.upload")+fileDir, 0777)
			}

			// Create file
			f, err = os.Create(viper.GetString("file.upload") + fileName)

			if err != nil {
				utils.LogError("Unable to create file")
				stream.SendAndClose(&protocol.UploadFileResponse{
					Message:          "Unable to create file ",
					UploadStatusCode: protocol.UploadStatusCode_Failed,
				})
				return
			}
			defer f.Close()

			firstChunk = false
		}

		// Write into file
		fmt.Println("name" + f.Name())
		err = utils.WriteToFile(f, chunks.Content)
		if err != nil {
			utils.LogError("Unable to write chunk of filename :" + err.Error())
			stream.SendAndClose(&protocol.UploadFileResponse{
				Message:          "Unable to write chunk of filename :",
				UploadStatusCode: protocol.UploadStatusCode_Failed,
			})
			return
		}
	}

	err = stream.SendAndClose(&protocol.UploadFileResponse{
		Message:          "Upload received with success",
		UploadStatusCode: protocol.UploadStatusCode_Ok,
	})

	if err != nil {
		err = errors.Wrapf(err,
			"Failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :")
	return
}

// DownloadFile handle file download call
func (c *CoreService) DownloadFile(params *protocol.DownloadFileParams, stream protocol.CoreService_DownloadFileServer) (err error) {
	var (
		writing = true
		buf     []byte
		n       int
		f       *os.File
	)

	fileFullPath := viper.GetString("file.upload") + params.FileName

	// Check if file exists and open
	f, err = os.Open(fileFullPath)
	defer f.Close()
	if err != nil {
		// File not found, send 404
		err = errors.Wrapf(err,
			"File does not exist")
		return
	}

	buf = make([]byte, 512)
	for writing {
		n, err = f.Read(buf)
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

		err = stream.Send(&protocol.ResChunk{
			Content: buf[:n],
		})
		if err != nil {
			err = errors.Wrapf(err,
				"Failed to send chunk via stream")
			return
		}
	}

	return
}
