package socker

import "github.com/pkg/errors"

type SockerClient struct {
	handlers []func(message []byte) bool
	index    int
}

func NewClient() SockerClient {
	return SockerClient{}
}

// Advance the state to the next handler
// panics if out of range.
func (c *SockerClient) Next() error {
	if c.index + 1 >= len(c.handlers) {
		return errors.New("Out of range exception")
	} else {
		c.index++
	}
	return nil
}

// Add a state to the client
// handler :
// a function that handles the []byte message
//		returns bool - true  if should advance to next handler
//					 - false if should not advance
func (c *SockerClient) Add(handler func(message []byte) bool) {
	c.handlers = append(c.handlers, handler)
}

// Handle the binary message
func (c *SockerClient) Handle(message []byte) error {
	if c.handlers[c.index](message) {
		return c.Next()
	}

	return nil
}
