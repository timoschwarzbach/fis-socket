package fissync

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
)

func Testsynclambda() {
	client, success := CreateMinioClient()
	if !success {
		log.Fatal("Failed to create Minio client")
	}
	// sync bucket "fis"
	downloadStaticFiles(client)
}

func downloadStaticFiles(client *minio.Client) {
	// Download the static files from the bucket
	// "fis" to the local filesystem
	// Define the bucket and local directory
	bucketName := "fis"
	localDir := "./static-test"

	// List objects in the bucket
	ctx := context.Background()
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}

		// Download each object
		localFilePath := filepath.Join(localDir, object.Key)
		err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}

		err = client.FGetObject(ctx, bucketName, object.Key, localFilePath, minio.GetObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Successfully downloaded %s\n", object.Key)
	}
}
