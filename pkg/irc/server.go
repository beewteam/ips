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
	fmt.Printf("Connecting to %s at port %s\n", s.Hostname, s.Port)

	conn, err := net.DialTimeout("tcp", s.Hostname+":"+s.Port, time.Second*10)
	if err != nil {
		fmt.Println(err)
		return false
	}

	s.Conn = conn

	return true
}

func (s *Server) String() string {
	return "hostname: " + s.Hostname + " port: " + s.Port
}
