package service

import (
	"dim-fs/protocol"
	"dim-fs/utils"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CoreService contains core file system service
type CoreService struct{}

//writeToFp takes in a file pointer and byte array and writes the byte array into the file
//returns error if pointer is nil or error in writing to file
func writeToFile(f *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {

		nw, err := f.Write(data[w:])
		if err != nil {
			return err
		}
		w += nw
		if nw >= n {
			return nil
		}
	}
}

// UploadFile handle upload file call
func (c *CoreService) UploadFile(stream protocol.CoreService_UploadFileServer) (err error) {
	firstChunk := true
	var f *os.File
	var chunks *protocol.Chunk

	for {

		chunks, err = stream.Recv() //ignoring the data  TO-Do save files received

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

		err = writeToFile(f, chunks.Content)
		if err != nil {
			utils.LogError("Unable to write chunk of filename :" + err.Error())
			stream.SendAndClose(&protocol.UploadFileResponse{
				Message:          "Unable to write chunk of filename :",
				UploadStatusCode: protocol.UploadStatusCode_Failed,
			})
			return
		}
	}

	//s.logger.Info().Msg("upload received")
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
