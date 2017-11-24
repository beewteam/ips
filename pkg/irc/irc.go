package irc

import (
	"fmt"
	"net"
)

func Pong(c net.Conn, pong string) bool {
	fmt.Fprintf(
		c,
		"PONG :%s\r\n",
		pong)
	return true
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
