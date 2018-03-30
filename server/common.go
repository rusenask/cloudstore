package server

import (
	"fmt"

	context "golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func getFileName(ctx context.Context) (filename string, err error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["filename"]) > 0 {
			return md["filename"][0], nil
		}
		return "", fmt.Errorf("filename not specified")
	}
	return "", fmt.Errorf("filename not specified")
}
