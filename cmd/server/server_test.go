package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetRoomFromName(t *testing.T) {
	assert := assert.New(t)
	var testRoom *Room
	var s Server
	assert.Equal(testRoom, s.getRoomFromName("test"), "Should be nil")
}
