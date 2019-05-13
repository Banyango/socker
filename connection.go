package socker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"log"
)

type SockerClientConnection struct {
	OpenHandler  func()
	CloseHandler  func()
	ReadHandlers []func(message []byte) bool
	index        int

	Connection *websocket.Conn

	inboundTCPChannel  chan []byte
	outboundTCPChannel chan []byte
	doneCh             chan bool
}

func NewClientConnection(connection *websocket.Conn) SockerClientConnection {
	return SockerClientConnection{
		Connection:connection,
		inboundTCPChannel:make(chan []byte),
		outboundTCPChannel:make(chan []byte),
		doneCh:make(chan bool),
	}
}

// Advance the state to the next handler
// returns err if out of range.
// Advance the state to the next handler
// panics if out of range.
func (c *SockerClientConnection) Next() error {
	if c.index + 1 >= len(c.ReadHandlers) {
		return errors.New("Out of range exception")
	} else {
		c.index++
	}
	return nil
}

// Add a state to the SockerClientConnection
// handler :
// a function that handles the []byte message
//		returns bool - true  if should advance to next handler
//					 - false if should not advance
func (c *SockerClientConnection) Add(handler func(message []byte) bool) {
	c.ReadHandlers = append(c.ReadHandlers, handler)
}

// Handle the binary message
func (c *SockerClientConnection) Handle() {
	select {
	case message, ok := <-c.inboundTCPChannel:
		if ok {
			if message != nil && c.ReadHandlers[c.index](message) {
				if err := c.Next(); err != nil {
					panic(err)
				}
			}
		}
	default:
	}
}

func (c *SockerClientConnection) Write(bytes []byte) {
	c.outboundTCPChannel <- bytes
}

func (c *SockerClientConnection) ReadPump() {
	defer func() {
		err := c.Connection.Close()
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}()
	for {
		select {
		case <-c.doneCh:
			c.doneCh <- true
			return

		default:
			messageType, data, err := c.Connection.ReadMessage()
			if err != nil {
				log.Println(err)

				c.doneCh <- true
			} else if messageType != websocket.BinaryMessage {
				log.Println("Non binary message recived, ignoring")
			}
			c.inboundTCPChannel <- data
		}
	}
}

func (c *SockerClientConnection) WritePump() {

	defer func() {
		fmt.Println("Connection Closing..")
		c.Connection.Close()
	}()
	for {
		select {
		case bytesOut := <-c.outboundTCPChannel:
			if err := c.Connection.WriteMessage(websocket.BinaryMessage, bytesOut); err != nil {
				log.Println(err)
			}
		case <-c.doneCh:
			c.doneCh <- true
			return
		}
	}
}







