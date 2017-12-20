import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
     server := NewServer()
     if assert.NotNil(t, server) {
          assert.Equal(t, server.Hostname), 7)
          assert.True(t, server.Port < 65535)
     }
}
