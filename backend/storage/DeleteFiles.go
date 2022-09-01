package storage

import (
	"context"
)

func DeleteFile(name string) error {

	ctx := context.Background()

	app, err := GetFirebaseApp()
	if err != nil {
		return err
	}

	bucket, err := GetBucket(app)
	if err != nil {
		return err
	}

	if err := bucket.Object(name).Delete(ctx); err != nil {
		return err
	}

	return nil
}
