package main

import (
	"fis/socket/ibis"
	"fis/socket/sequences"
	"fis/socket/socket"
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
	fmt.Println("socket process started")

	go func() {
		fmt.Println("Starting file server")
		// http.ListenAndServe(":8080", http.FileServer(http.Dir("static")))
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
		}))
	}()

	// start socket server
	server := socket.StartSocket()

	// create ibis controller
	go func() {
		fmt.Println("creating ibis controller")
		ibis.CreateController(server)
	}()

	// create media controller
	fmt.Println("creating media controller")
	sequences.CreateController(server)

	fmt.Println("socket process ended")
}
