package storage

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/iam"
	gcstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
	iambp "google.golang.org/genproto/googleapis/iam/v1"
)

var err error

var app *firebase.App
var appErr error
var onceApp sync.Once

func GetFirebaseApp() (*firebase.App, error) {
	app = nil

	onceApp.Do(func() {
		sa := option.WithCredentialsFile("./posthis-go-firebase-adminsdk-h1r8j-14cc5fa9c9.json")
		app, err = firebase.NewApp(context.Background(), nil, sa)
		if err != nil {
			appErr = err
			return
		}
	})

	return app, appErr
}

var bucket *gcstorage.BucketHandle
var bucketErr error
var onceBucket sync.Once

func GetBucket(app *firebase.App) (*gcstorage.BucketHandle, error) {

	onceBucket.Do(func() {

		ctx := context.Background()

		client, err := app.Storage(ctx)
		if err != nil {
			bucketErr = err
			return
		}

		bucket, bucketErr = client.Bucket(os.Getenv("POSTHIS_APP_STORAGE_BUCKET_NAME"))
	})

	return bucket, bucketErr
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

func Init() error {

	app, err := GetFirebaseApp()
	if err != nil {
		return err
	}

	bucket, err := GetBucket(app)
	if err != nil {
		return err
	}

	if err = makeBucketPublic(bucket); err != nil {
		return err
	}

	return nil
}
