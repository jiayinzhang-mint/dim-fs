package service

import (
	"dim-fs/protocol"
	"dim-fs/utils"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

// CoreService contains core file system service
type CoreService struct{}

//writeToFp takes in a file pointer and byte array and writes the byte array into the file
//returns error if pointer is nil or error in writing to file
func writeToFp(fp *os.File, data []byte) error {
	w := 0
	n := len(data)
	for {

		nw, err := fp.Write(data[w:])
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
	var fp *os.File

	var fileData *protocol.Chunk

	for {

		fileData, err = stream.Recv() //ignoring the data  TO-Do save files received

		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}

		if firstChunk { //first chunk contains file name

			fp, err = os.Create("output")

			if err != nil {
				utils.LogError("Unable to create file")
				stream.SendAndClose(&protocol.UploadFileResponse{
					Message:          "Unable to create file ",
					UploadStatusCode: protocol.UploadStatusCode_Failed,
				})
				return
			}
			defer fp.Close()

			firstChunk = false
		}

		err = writeToFp(fp, fileData.Content)
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
			"failed to send status code")
		return
	}
	fmt.Println("Successfully received and stored the file :")
	return
}
