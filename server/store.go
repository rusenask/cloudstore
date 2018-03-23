package server

import (
	"fmt"
	"io"

	"github.com/rusenask/cloudstore"
	context "golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

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

func getFileName(ctx context.Context) (filename string, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["filename"]) > 0 {
			return md["filename"][0], nil
		}
		return "", fmt.Errorf("filename not specified")
	}
	return "", fmt.Errorf("filename not specified")
}
