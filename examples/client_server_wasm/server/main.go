package main

import (
	"context"
	"fmt"
	"github.com/Banyango/socker"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
)

func connectionMiddleware(next http.Handler, server *Server) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), "Server", server)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
	return true
}}

func websocketFunc(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil	)
	if err != nil {
		log.Println(err)
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		log.Println(code, text)
		return nil
	})

	server := r.Context().Value("Server").(*Server)

	connection := socker.NewClientConnection(conn)

	connection.Add(func(message []byte) bool {

		fmt.Println("One", string(message))

		connection.Write([]byte("hi"))

		if err != nil {
			panic("Something bad happened!")
		}

		return true
	})

	connection.Add(func(message []byte) bool {

		fmt.Println("Three", string(message))

		connection.Write([]byte("hi again"))

		if err != nil {
			panic("Something bad happened!")
		}

		os.Exit(0)

		return false

	})

	server.Register <- connection

	go connection.ReadPump()
	go connection.WritePump()

}

type Server struct {
	Clients []*socker.SockerClientConnection
	Register chan *socker.SockerClientConnection
	mux sync.Mutex
}

func main() {
	server := Server{}

	server.Register = make(chan *socker.SockerClientConnection)

	http.Handle("/ws", connectionMiddleware(http.HandlerFunc(websocketFunc), &server))

	//http.HandleFunc("/", websocketFunc)

	go mainLoop(&server)

	log.Fatal(http.ListenAndServe(":3000",nil))
}

// this would be the main processing loop of the game.
func mainLoop(server *Server) {
	for {

		for _, value := range server.Clients {
			value.Handle()
		}

		select {
		case client, ok := <- server.Register:
			if ok {
				server.Clients = append(server.Clients, client)
			}
		default:
		}
	}
}