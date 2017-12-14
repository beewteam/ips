package irc

import (
	"fmt"
	"io"
	"net"
	"strings"
)

type ircCommand struct {
	name     string
	shortcut string
	handler  func(c *Client, msg string) bool
}

var (
	ircCommands []ircCommand
)

func isCommand(name string, cmd *ircCommand) bool {
	return cmd.name == name || cmd.shortcut == name
}

func SendCommand(w io.Writer, msg string) {
	reply := msg + "\r\n"
	fmt.Fprint(w, reply)
}

func PingHandler(c *Client, msg string) bool {
	words := strings.Fields(msg)

	if len(words) < 1 {
		return false
	}

	reply := "PONG" + " " + words[1]
	if len(words) == 3 {
		reply = reply + words[2]
	}

	SendCommand(c.Server.Conn, reply)

	return true
}

func Handle(c *Client, msg string) {
	words := strings.Fields(msg)
	for _, cmd := range ircCommands {
		if isCommand(words[0], &cmd) {
			cmd.handler(c, msg)
		}
	}
}

func Init() {
	ircCommands = []ircCommand{
		{"PING", "PG", PingHandler},
	}
}

func Reg(c net.Conn, nick string, username string, fullname string) bool {
	fmt.Fprintf(
		c,
		"NICK %s\r\nUSER %s 0 0 :%s\r\n",
		nick, username, fullname)
	return true
}

func Join(c net.Conn, channel string) bool {
	fmt.Fprintf(
		c,
		"JOIN %s\r\n",
		channel)
	return true
}

func Part(c net.Conn, data string) bool {
	fmt.Fprintf(
		c,
		"PART %s\r\n",
		data)
	return true
}

func Nick(c net.Conn, nick string) bool {
	fmt.Fprintf(
		c,
		"NICK %s\r\n",
		nick)
	return true
}

func Quit(c net.Conn, quitMsg string) bool {
	fmt.Fprintf(
		c,
		"QUIT :%s\r\n",
		quitMsg)
	return true
}

func Topic(c net.Conn, channel string, data string) bool {
	fmt.Fprintf(
		c,
		"TOPIC %s :%s\r\n",
		channel, data)
	return true
}

func Action(c net.Conn, channel string, data string) bool {
	fmt.Fprintf(
		c,
		"PRIVMSG %s :\001ACTION %s\001\r\n",
		channel, data)
	return true
}

func Msg(c net.Conn, channel string, data string) bool {
	fmt.Fprintf(
		c,
		"PRIVMSG %s :%s\r\n",
		channel, data)
	return true
}
