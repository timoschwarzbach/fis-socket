package socket

import (
	"fmt"
	"net/http"

	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

type Server = socket.Server

func StartSocket() *Server {
	io := socket.NewServer(nil, nil)
	io.Opts().SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	http.Handle("/socket.io/", io.ServeHandler(nil))
	go func() {
		if err := http.ListenAndServe(":3000", nil); err != nil {
			fmt.Printf("ListenAndServe: %v\n", err)
		}
	}()
	fmt.Println("Starting socket.io server on :3000")

	io.On("connection", func(clients ...any) {
		fmt.Println("New connection")
		client := clients[0].(*socket.Socket)
		client.On("event", func(datas ...any) {
		})
		client.On("disconnect", func(...any) {
			fmt.Println("A client disconnected")
		})
	})

	return io
}
