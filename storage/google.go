package storage

import (
	"context"
	"errors"
	"io"
	"os"

	"google.golang.org/api/option"

	gstorage "cloud.google.com/go/storage"

	log "github.com/sirupsen/logrus"
)

type GoogleCloudStorage struct {
	Storage
	Options    []option.ClientOption
	client     *gstorage.Client
	bucket     *gstorage.BucketHandle
	bucketName string
}

func (gcs *GoogleCloudStorage) Setup() error {
	bktName := os.Getenv("GOOGLE_STORAGE_BUCKET")
	if bktName == "" {
		return errors.New("GOOGLE_STORAGE_BUCKET env must be set")
	}

	gcs.bucketName = bktName

	projectID := os.Getenv("GOOGLE_STORAGE_PROJECT_ID")
	if projectID == "" {
		return errors.New("GOOGLE_STORAGE_PROJECT_ID env must be set")
	}

	location := os.Getenv("GOOGLE_STORAGE_LOCATION")
	if location == "" {
		return errors.New("GOOGLE_STORAGE_LOCATION env must be set")
	}

	ctx := context.Background()
	client, err := gstorage.NewClient(ctx, gcs.Options...)
	if err != nil {
		return err
	}

	gcs.client = client

	bkt := client.Bucket(bktName)

	attrs, err := bkt.Attrs(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"bucket_name": bktName,
		}).Error("failed to access bucket")
		return err
	}
	log.Infof("bucket %s, created at %s, is located in %s with storage class %s\n",
		attrs.Name, attrs.Created, attrs.Location, attrs.StorageClass)

	gcs.bucket = bkt

	return nil
}

func (gcs *GoogleCloudStorage) PublicURL(filename string) string {
	return "https://storage.googleapis.com/" + gcs.bucketName + "/" + filename
}

func (gcs *GoogleCloudStorage) Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error {
	o := gcs.bucket.Object(filename)
	w := o.NewWriter(ctx)

	w.ObjectAttrs = gstorage.ObjectAttrs{
		Name:     filename,
		Metadata: metadata,
	}

	_, err := w.Write(data)
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) Get(ctx context.Context, filename string) (io.ReadCloser, error) {
	o := gcs.bucket.Object(filename)
	return o.NewReader(ctx)
}

func (gcs *GoogleCloudStorage) Delete(ctx context.Context, filename string) error {
	o := gcs.bucket.Object(filename)
	return o.Delete(ctx)
}
