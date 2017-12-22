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

func TestGetClientFromName(t *testing.T) {
	assert := assert.New(t)
	var testClient1 *Client
	var mockServer = Server{
		make(map[*Client]bool),
		make(map[*Room]bool),
		sync.Mutex{},
	}
	assert.Equal(testClient1, mockServer.getClientFromName("test"), "Should be nil")

	var testClient2 = Client{
		"Mad Client",
		make(map[*Room]bool),
	}
	mockServer.clients[&testClient2] = true
	assert.Equal(&testClient2, mockServer.getClientFromName("Mad Client"), "Should be not nil")
}
