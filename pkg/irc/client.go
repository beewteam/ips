package irc

import (
	"bufio"
	"fmt"
	"time"
)

type Client struct {
	Username string
	FullName string
	Channel  string
	Nick     string
	Data     string
	Server   Server
}

func (c *Client) Login(nick string) bool {
	return Reg(c.Server.Conn, nick, c.Username, c.FullName)
}

func (c *Client) JoinChannel(channel string) bool {
	c.Channel = channel
	return Join(c.Server.Conn, channel)
}

func (c *Client) LeaveChannel() bool {
	return Part(c.Server.Conn, c.Channel)
}

func (c *Client) HandleData() bool {
	message, _ := bufio.NewReader(c.Server.Conn).ReadString('\r')
	if len(message) > 0 {
		fmt.Printf("Message:%s\n", message)
	}
	return true
}

func (c *Client) LogMessage(nick string, msg string) bool {
	fmt.Printf(
		"%s - [%s] <%s> %s\n",
		time.Now().String(), c.Channel, c.Nick, msg)
	return true
}

func (c *Client) Close() bool {
	c.Server.Conn.Close()
	return true
}

func (c *Client) Auth(passwd string) bool {
	return Msg(c.Server.Conn, "NickServ", "IDENTIFY "+passwd)
}
