package main

import (
	"fmt"
	"github.com/Banyango/socker"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/ws", nil)

	if err != nil {
		log.Fatal(err)
	}

	client := socker.NewClient()

	client.Add(func(message []byte) bool {
		fmt.Println("two")
		err = conn.WriteMessage(websocket.BinaryMessage, []byte("hi"))
		if err != nil {
			log.Fatal(err)
		}
		return true
	})

	client.Add(func(message []byte) bool {
		fmt.Println("done")
		os.Exit(0)
		return false
	})

	err = conn.WriteMessage(websocket.BinaryMessage, []byte("hi"))

	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		default:
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			err = client.Handle(data)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}


