package storage

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var once sync.Once
var bh *storage.BucketHandle
var err error

func GetBucketHandle() (*storage.BucketHandle, error) {

	once.Do(func() {

		bh = nil

		config := &firebase.Config{
			StorageBucket: "posthis-f7302.appspot.com",
		}

		opt := option.WithCredentialsFile("posthis-f7302-firebase-adminsdk-oha7v-cd28dd043b.json")
		app, err := firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			log.Fatalln(err)
			return
		}

		client, err := app.Storage(context.Background())
		if err != nil {
			log.Fatalln(err)
			return
		}

		bh, err = client.DefaultBucket()
		if err != nil {
			log.Fatalln(err)
			return
		}
	})

	return bh, err
}
