package utils

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"posthis/entity"
	"posthis/storage"
	"time"

	gcstorage "cloud.google.com/go/storage"

	uuid "github.com/satori/go.uuid"
)

type Media = entity.Media

func UploadMultipleFiles(files []*multipart.FileHeader, media *[]*Media) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	bucket, err := storage.GetBucketHandle()
	if err != nil {
		return err
	}

	for i := range files {

		file, err := files[i].Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_bytes := make([]byte, files[i].Size)

		_, err = file.Read(_bytes)
		if err != nil {
			return err
		}

		mime := http.DetectContentType(_bytes)

		//random name
		name := uuid.NewV4().String() + files[i].Filename
		object := bucket.Object(name)
		writer := object.NewWriter(ctx)
		defer writer.Close()

		if _, err := io.Copy(writer, bytes.NewReader(_bytes)); err != nil {
			return err
		}

		if err := object.ACL().Set(context.Background(), gcstorage.AllUsers, gcstorage.RoleReader); err != nil {
			return err
		}

		*media = append(*media, &Media{Name: name, Mime: mime})
	}

	return nil
}
