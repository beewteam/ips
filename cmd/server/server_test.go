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
	var server = Server{
		make(map[*Client]bool),
		make(map[*Room]bool),
		sync.Mutex{},
	}
	assert.Equal(testRoom1, server.getRoomFromName("test"), "Should be nil")

	var testRoom2 = Room{
		"cool room",
		make(map[*Client]bool),
	}
	server.rooms[&testRoom2] = true
	assert.Equal(&testRoom2, server.getRoomFromName("cool room"), "Should be not nil")
}

func TestInvalidResponse(t *testing.T) {
	mock := connMockObject{}

	testServer := Server{}
	testString := []string{"ol", "mol", "tol", "rol", "poi"}
	command := Command{}
	command.client = &Client{}
	command.client.conn = &mock

	mock.On("Write", []byte("Invalid syntax\n")).Return(15, nil)
	handleUserCommand(&testServer, testString, &command)

	mock.AssertExpectations(t)
}

func TestClientGetSet(t *testing.T) {
	assert := assert.New(t)
	client := Client{}

	assert.Equal(false, setClientNickname(&client, ""), "Should be false")

	assert.Equal("", getClientNickname(nil), "Should be false")
}

func TestNickCommand(t *testing.T) {
	mock := connMockObject{}

	testServer := Server{}
	testString := []string{"ol", "mol"}
	command := Command{}
	command.client = &Client{}
	command.client.conn = &mock

	mock.On("Write", []byte("Invalid syntax\n")).Return(15, nil)
	handleUserCommand(&testServer, testString, &command)

	mock.AssertExpectations(t)
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
	if err12 != nil {
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
