package fissync

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
)

type SyncController struct {
	MinioClient *minio.Client
	interval    int
	dbSync      chan bool
}

func CreateSynchronizer(dbSync chan bool) *SyncController {
	client, success := CreateMinioClient()
	if !success {
		log.Fatalln("SyncService:\tFailed to create Minio client")
	}

	log.Println("SyncService:\tCreating SyncService")
	return &SyncController{
		MinioClient: client,
		interval:    get_sync_interval(),
		dbSync:      dbSync,
	}
}

// retrieves the synchronization interval from the environment
// if the environment variable is not set or cannot be parsed,
// it returns a default interval of 60 seconds.
func get_sync_interval() int {
	env_interval, exists := os.LookupEnv("SYNC_INTERVAL_SECONDS")
	if !exists {
		log.Println("SyncService:\tSetup:\tenvironment variable SYNC_INTERVAL_SECONDS not specified. using default")
		return 60
	}

	interval, err := strconv.Atoi(env_interval)
	if err != nil {
		log.Println("SyncService:\tSetup:\tFailed to parse SYNC_INTERVAL_SECONDS. using default")
		return 60
	}

	return interval
}

func (s *SyncController) StartIntervalBackgroundSync() {
	log.Println("SyncService:\tStarting SyncService")
	go func() {
		ticker := time.NewTicker(time.Duration(s.interval) * time.Second)
		defer ticker.Stop()

		s.Sync()
		for range ticker.C {
			log.Println("SyncService:\tSyncing to upstream...")
			s.Sync()
		}
	}()
}

func (s *SyncController) Sync() {
	s.syncDatabase()
	s.syncStaticFiles("fis", "./static")
	s.syncStaticFiles("tagesschau", "./static/tagesschau")
}

func (s *SyncController) syncDatabase() bool {
	endpoint, exists := os.LookupEnv("MANAGE_ENDPOINT")
	if !exists {
		log.Fatalln("SyncService:\tfis-manage endpoint not specified")
	}

	log.Println("SyncService:\tDownloading database")
	err := DownloadFile("database.sqlite", endpoint)
	if err != nil {
		log.Println("SyncService:\tFailed to download database")
		log.Println(err)
		return false
	}

	// notify the database service to reload the database
	log.Println("SyncService:\tTelling database service to reload the database")
	go func() {
		s.dbSync <- true
	}()

	return true
}

func (s *SyncController) syncStaticFiles(bucketName string, localDir string) {
	// Sync the static files from the bucket
	// "fis" to the local filesystem
	// Define the bucket and local directory

	// List objects in the bucket
	ctx := context.Background()
	objectCh := s.MinioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

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
			err = s.MinioClient.FGetObject(ctx, bucketName, object.Key, localFilePath, minio.GetObjectOptions{})
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("SyncService:\tResource:\tSuccessfully synced %s\n", object.Key)
		} else {
			// log.Printf("SyncService:\tResource:\tAlready up-to-date: %s\n", object.Key)
		}

		// Add the object to the map
		bucketObjects[localFilePath] = struct{}{}
	}

	// Remove local files that are not in the bucket
	err := filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if _, exists := bucketObjects[path]; !exists {
				err := os.Remove(path)
				if err != nil {
					log.Fatalln(err)
				}
				// todo: only remove if sqlite sync was successful
				log.Printf("SyncService:\tResource:\tRemoved local file: %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}
