package service

import (
	"dim-fs/protocol"
	"dim-fs/utils"
	"io"

	"github.com/pkg/errors"
)

// CoreService contains core file system service
type CoreService struct{}

// UploadFile handle upload file call
func (c *CoreService) UploadFile(stream protocol.CoreService_UploadFileServer) (err error) {
	for {
		_, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			err = errors.Wrapf(err,
				"failed unexpectadely while reading chunks from stream")
			return
		}
		utils.LogInfo("upload received")
	}

END:
	err = stream.SendAndClose(&protocol.UploadFileResponse{
		Message:          "Upload received with success",
		UploadStatusCode: protocol.UploadStatusCode_Ok,
	})
	if err != nil {
		err = errors.Wrapf(err,
			"failed to send status code")
		return
	}

	return
}
