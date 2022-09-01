package storage

import (
	"context"
	"io"
	"mime/multipart"
	"posthis/viewmodel"
	"time"

	uuid "github.com/satori/go.uuid"
)

type MediaCreateVM = viewmodel.MediaCreateVM

func UploadMultipleFiles(files []*multipart.FileHeader) ([]MediaCreateVM, error) {
	var media []MediaCreateVM

	if len(files) == 0 {
		return media, nil
	}

	// set timeout for file upload
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	app, err := GetFirebaseApp()
	if err != nil {
		return nil, err
	}

	bucket, err := GetBucket(app)
	if err != nil {
		return nil, err
	}

	for _, ptrFile := range files {
		file, err := ptrFile.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// unique name
		Name := uuid.NewV4().String() + ptrFile.Filename
		object := bucket.Object(Name)
		writer := object.NewWriter(ctx)

		// Copy file to storage
		_, err = io.Copy(writer, file)
		if err != nil {
			return nil, err
		}

		if err := writer.Close(); err != nil {
			return nil, err
		}

		// get attributes
		attrs, err := object.Attrs(ctx)
		if err != nil {
			return nil, err
		}

		Mime := attrs.ContentType
		Url := attrs.MediaLink

		// create media with its respective characteristics
		media = append(media, MediaCreateVM{Name, Mime, Url})
	}

	return media, nil
}
