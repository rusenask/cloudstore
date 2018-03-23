package server

import (
	"github.com/rusenask/cloudstore"
	context "golang.org/x/net/context"
)

func (s CloudStorageServiceServer) Delete(ctx context.Context, r *cloudstore.DeleteRequest) (*cloudstore.DeleteResponse, error) {
	// return nil, errors.New("not yet implemented")
	err := s.storage.Delete(ctx, r.Filename)
	if err != nil {
		return nil, err
	}

	return &cloudstore.DeleteResponse{
		Filename: r.Filename,
	}, nil
}
