package main

import (
	"fis/socket/ibis"
	"fis/socket/sequences"
	"fis/socket/socket"
	fissync "fis/socket/sync"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	log.Println("General:\tsocket process started")

	go func() {
		fmt.Println("FileServer:\tStarting file server")
		// http.ListenAndServe(":8080", http.FileServer(http.Dir("static")))
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
		}))
	}()

	// start socket server
	server := socket.StartSocket()

	// start sync service
	go func() {
		sync := fissync.CreateSynchronizer()
		sync.StartIntervalBackgroundSync()
	}()

	// create ibis controller
	go func() {
		fmt.Println("IbisService:\tcreating ibis controller")
		ibis.CreateController(server)
	}()

	// create media controller
	fmt.Println("SequenceController:\tcreating media controller")
	sequences.CreateController(server)

	fmt.Println("socket process ended")
}
