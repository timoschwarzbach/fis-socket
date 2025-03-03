package main

import (
	"fis/socket/content"
	"fis/socket/ibis"
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

	// channel for database sync messages
	dbSync := make(chan bool, 1)

	// start sync service
	sync := fissync.CreateSynchronizer(dbSync)
	sync.StartBackgroundSync()

	// create ibis controller
	go func() {
		fmt.Println("IbisService:\tcreating ibis controller")
		ibis.CreateController(server)
	}()

	// create media controller
	fmt.Println("SequenceController:\tcreating media controller")
	content.CreateController(server, dbSync)

	fmt.Println("socket process ended")
}
