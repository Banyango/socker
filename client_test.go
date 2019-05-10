package socker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_Handle(t *testing.T) {
	client := NewClient()

	counter := 0
	client.Add(func(message []byte) bool {
		counter = 1
		return false
	})

	err := client.Handle(nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, counter)
}

func TestClient_HandleNext(t *testing.T) {
	client := NewClient()

	counter := 0
	client.Add(func(message []byte) bool {
		counter = 1
		return true
	})

	client.Add(func(message []byte) bool {
		counter = 2
		return false
	})


	err := client.Handle(nil)

	assert.NoError(t, err)

	err = client.Handle(nil)

	assert.NoError(t, err)
	assert.Equal(t, 2, counter)
}

func TestClient_HandleReturnFalse(t *testing.T) {
	client := NewClient()

	counter := 0
	client.Add(func(message []byte) bool {
		counter = 1
		return false
	})

	client.Add(func(message []byte) bool {
		counter = 2
		return false
	})


	err := client.Handle(nil)

	assert.NoError(t, err)

	err = client.Handle(nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, counter)
}

func TestClient_HandleReturnIndexOutOfBounds(t *testing.T) {
	client := NewClient()

	counter := 0
	client.Add(func(message []byte) bool {
		counter = 1
		return true
	})

	err := client.Handle(nil)

	assert.Equal(t, 1, counter)
	assert.Error(t, err)
}

func TestClient_HandleAppend(t *testing.T) {
	client := NewClient()

	client.Add(func(message []byte) bool {
		return true
	})

	assert.Equal(t, 1, len(client.handlers))
}