package fissync

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Panicln("DownloadFile:\tCannot create file")
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		log.Println(err.Error())
		log.Println(url)
		log.Panicln("DownloadFile:\tCannot get file from url")
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panicln("DownloadFile:\tCannot write to file")
		return err
	}

	return nil
}
