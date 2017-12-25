package main

import (
	"sync"
	"testing"

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
