package server

import (
	"io"

	"github.com/rusenask/cloudstore"

	log "github.com/sirupsen/logrus"
)

func (s CloudStorageServiceServer) Store(stream cloudstore.CloudStorageService_StoreServer) error {
	filename, err := getFileName(stream.Context())
	if err != nil {
		return err
	}

	// receiving file
	var blob []byte
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				goto END
			}

			log.Errorf("failed while reading chunks from stream: %s", err)
			return err
		}
		blob = append(blob, req.Content...)
	}

END:
	// storing
	err = s.storage.Store(
		stream.Context(),
		filename,
		blob,
		map[string]string{},
	)

	if err != nil {
		return err
	}

	return stream.SendAndClose(&cloudstore.UploadResponse{
		Message: "Upload received with success",
		Url:     s.storage.PublicURL(filename),
		Code:    cloudstore.UploadStatusCode_ok,
	})
}
