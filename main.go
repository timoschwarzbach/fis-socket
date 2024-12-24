package main

import (
	"fis/socket/sequences"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")

	go func() {
		fmt.Println("Starting file server")
		// http.ListenAndServe(":8080", http.FileServer(http.Dir("static")))
		http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			http.FileServer(http.Dir("static")).ServeHTTP(w, r)
		}))
	}()
	sequences.CreateController()

	fmt.Println("Bye, World!")
}
