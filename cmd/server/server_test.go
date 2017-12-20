package main

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRoomFromName(t *testing.T) {
	assert := assert.New(t)
	var testRoom1 *Room
	var mockServer = Server{
		make(map[*Client]bool),
		make(map[*Room]bool),
		sync.Mutex{},
	}
	assert.Equal(testRoom1, mockServer.getRoomFromName("test"), "Should be nil")

	var testRoom2 = Room{
		"cool room",
		make(map[*Client]bool),
	}
	mockServer.rooms[&testRoom2] = true
	assert.Equal(&testRoom2, mockServer.getRoomFromName("cool room"), "Should be not nil")
}
