package fissync

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateMinioClient() (*minio.Client, bool) {
	endpoint, exists := os.LookupEnv("MINIO_ENDPOINT")
	if !exists {
		log.Fatalln("SyncService:\tSetup:\tminio endpoint not specified")
		return nil, false
	}

	accessKeyID, exists := os.LookupEnv("MINIO_ACCESS_KEY")
	if !exists {
		log.Fatalln("SyncService:\tSetup:\tminio access key not specified")
		return nil, false
	}

	secretAccessKey, exists := os.LookupEnv("MINIO_SECRET_KEY")
	if !exists {
		log.Fatalln("SyncService:\tSetup:\tminio secret key not specified")
		return nil, false
	}

	useSSLString, exists := os.LookupEnv("MINIO_SSL")
	if !exists {
		log.Fatalln("SyncService:\tSetup:\tminio secret key not specified")
		return nil, false
	}
	useSSL := useSSLString == "true"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("SyncService:\tSetup:\tminioClient initialized successfully")
	return minioClient, true
}
