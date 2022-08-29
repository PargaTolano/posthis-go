package utils

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"posthis/entity"
	"time"

	"cloud.google.com/go/iam"
	gcstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	uuid "github.com/satori/go.uuid"

	"google.golang.org/api/option"
	iambp "google.golang.org/genproto/googleapis/iam/v1"
)

type Media = entity.Media

var app *firebase.App

func Init() error {
	var err error
	sa := option.WithCredentialsFile("./posthis-go-firebase-adminsdk-h1r8j-14cc5fa9c9.json")
	app, err = firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		return err
	}

	bucket, err := getBucket()
	if err != nil {
		return err
	}

	makeBucketPublic(bucket)

	return nil
}

func makeBucketPublic(bucket *gcstorage.BucketHandle) error {
	ctx := context.Background()

	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}

	role := "roles/storage.objectViewer"
	policy.Bindings = append(policy.Bindings, &iambp.Binding{
		Role:    role,
		Members: []string{iam.AllUsers},
	})

	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}

	return nil
}

func getBucket() (*gcstorage.BucketHandle, error) {
	ctx := context.Background()
	client, err := app.Storage(ctx)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(os.Getenv("POSTHIS_APP_STORAGE_BUCKET_NAME"))
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func UploadMultipleFiles(files []*multipart.FileHeader, media *[]*Media) error {
	// set timeout for file upload
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	// get storage bucket
	bucket, err := getBucket()
	if err != nil {
		return err
	}
	for _, ptrFile := range files {
		file, err := ptrFile.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// write image to bucket
		name := uuid.NewV4().String() + ptrFile.Filename
		object := bucket.Object(name)
		writer := object.NewWriter(ctx)

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		if err := writer.Close(); err != nil {
			return err
		}

		// get attributes
		attrs, err := object.Attrs(ctx)
		if err != nil {
			return err
		}

		mime := attrs.ContentType
		url := attrs.MediaLink

		// create media with its respective characteristics
		*media = append(*media, &Media{Name: name, Mime: mime, Url: url})
	}

	return nil
}
