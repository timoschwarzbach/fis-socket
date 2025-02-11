package fissync

import (
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func CreateMinioClient() (*minio.Client, bool) {
	endpoint, exists := os.LookupEnv("MINIO_ENDPOINT")
	if !exists {
		fmt.Println("minio endpoint not specified")
		return nil, false
	}

	accessKeyID, exists := os.LookupEnv("MINIO_ACCESS_KEY")
	if !exists {
		fmt.Println("minio access key not specified")
		return nil, false
	}

	secretAccessKey, exists := os.LookupEnv("MINIO_SECRET_KEY")
	if !exists {
		fmt.Println("minio secret key not specified")
		return nil, false
	}

	useSSLString, exists := os.LookupEnv("MINIO_SSL")
	if !exists {
		fmt.Println("minio secret key not specified")
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

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	return minioClient, true
}
