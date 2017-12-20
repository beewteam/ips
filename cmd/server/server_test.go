package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetRoomFromName(t *testing.T) {
	assert := assert.New(t)
	var testRoom *Room
	assert.Equal(testRoom, getRoomFromName("test"), "Should be nil")
}
