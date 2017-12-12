package irc

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Hostname string
	Port     string
	Conn     net.Conn
}

func (s *Server) Connect() bool {
	conn, err := net.DialTimeout("tcp", s.Hostname+":"+s.Port, time.Second*10)
	if err != nil {
		fmt.Println(err)
		return false
	}

	s.Conn = conn

	return true
}
