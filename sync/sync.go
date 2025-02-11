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
	syncStaticFiles(client)
}

func syncStaticFiles(client *minio.Client) {
	// Sync the static files from the bucket
	// "fis" to the local filesystem
	// Define the bucket and local directory
	bucketName := "fis"
	localDir := "./static-test"

	// List objects in the bucket
	ctx := context.Background()
	objectCh := client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

	// Track the objects in the bucket
	bucketObjects := make(map[string]struct{})

	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}

		// Define the local file path
		localFilePath := filepath.Join(localDir, object.Key)

		// Check if the local file exists and is up-to-date
		shouldDownload := true
		if fileInfo, err := os.Stat(localFilePath); err == nil {
			// Compare the modification time
			objectModTime := object.LastModified
			localModTime := fileInfo.ModTime()
			if localModTime.After(objectModTime) || localModTime.Equal(objectModTime) {
				shouldDownload = false
			}
		}

		if shouldDownload {
			// Create the local directory if it doesn't exist
			err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm)
			if err != nil {
				log.Fatalln(err)
			}

			// Download the object
			err = client.FGetObject(ctx, bucketName, object.Key, localFilePath, minio.GetObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Successfully synced %s\n", object.Key)
		} else {
			fmt.Printf("Already up-to-date: %s\n", object.Key)
		}

		// Add the object to the map
		bucketObjects[localFilePath] = struct{}{}
	}

	// Remove local files that are not in the bucket
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if _, exists := bucketObjects[path]; !exists {
				err := os.Remove(path)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("Removed local file: %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
