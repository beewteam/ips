package main

import (
	"sync"
	"testing"
	"net"

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
	conn12, err12 := net.Dial("tcp", "golang.org:80")
	if err != nil {
	// handle error
	}
	var testClient2 = Client{
		conn12,
		make(map[*Room]bool),
		"Mad Client",
		"Mad Client",
		"Mad Client",
		"Host",		
	}
	mockServer.clients[&testClient2] = true
	assert.Equal(&testClient2, mockServer.getClientFromName("Mad Client"), "Should be not nil")
}
