package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Client struct {
	conn     net.Conn
	username string
	fullName string
	channel  string
	nick     string
	data     string
}

func (c *Client) Connect(server string, port string) bool {
	conn, err := net.DialTimeout("tcp", server+":"+port, time.Second*10)
	if err != nil {
		fmt.Println(err)
		return false
	}

	c.conn = conn

	return true
}

func (c *Client) Login(nick string) bool {
	return Reg(c.conn, nick, c.username, c.fullName)
}

func (c *Client) JoinChannel(channel string) bool {
	c.channel = channel
	return Join(c.conn, channel)
}

func (c *Client) LeaveChannel() bool {
	return Part(c.conn, c.channel)
}

func (c *Client) HandleData() bool {
	message, _ := bufio.NewReader(c.conn).ReadString('\r')
	fmt.Printf("Message:%s\n", message)

	return true
}

func (c *Client) LogMessage(nick string, msg string) bool {
	fmt.Printf(
		"%s - [%s] <%s> %s\n",
		time.Now().String(), c.channel, c.nick, msg)
	return true
}

func (c *Client) Close() bool {
	c.conn.Close()
	return true
}
