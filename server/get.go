package server

import (
	"io"

	"github.com/rusenask/cloudstore"

	log "github.com/sirupsen/logrus"
)

func (s CloudStorageServiceServer) Get(r *cloudstore.GetRequest, stream cloudstore.CloudStorageService_GetServer) error {

	// res := &cloudstore.Chunk{}
	ctx := stream.Context()

	// getting
	reader, err := s.storage.Get(ctx, r.Filename)
	if err != nil {
		log.Errorf("failed to retrieve file from storage: %s", err)
		return err
	}
	defer reader.Close()

	buf := make([]byte, 64*1024)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			v, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				log.WithFields(log.Fields{
					"buf":  string(buf),
					"err":  err,
					"read": v,
				}).Error("got error while reading from buf")
				return err
			}

			if v == 0 {
				return nil
			}
			//in case it is a string file, you could check its content here...
			if err := stream.Send(&cloudstore.Chunk{
				Content: buf[:v],
			}); err != nil {
				return err
			}
		}
	}
}
