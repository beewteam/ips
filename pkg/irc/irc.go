package irc

import (
	"fmt"
	"net"
	"strings"
)

type ircCommand struct {
	format string
}

type message struct {
}

type Communicator struct {
	socket        net.Conn
	processingMsg message
	messagesQueue chan string
	errors        chan string
}

const (
	t_cmd_name  = "@@CmdName@@"
	proto_delim = "\r\n"

	warning_error = "1"
)

var (
	ircCommands map[string]ircCommand
)

func (c *Communicator) Init() {
	ircCommands = map[string]ircCommand{
		"PING":    {t_cmd_name + " %s"},
		"PART":    {t_cmd_name + " %s"},
		"NICK":    {t_cmd_name + " %s"},
		"QUIT":    {t_cmd_name + " :%s"},
		"TOPIC":   {t_cmd_name + " %s :%s"},
		"PRIVMSG": {t_cmd_name + " %s :%s"},
	}
	c.messagesQueue = make(chan string)
	c.errors = make(chan string)
}

func wrapMessage(format string, cmdName string, params ...interface{}) string {
	parsedFormat := strings.Replace(format, t_cmd_name, cmdName, 1) + proto_delim
	return fmt.Sprintf(parsedFormat, params...)
}

func (c *Communicator) handleSendMessage() {
	for {
		if c.socket == nil {
			c.errors <- warning_error + " : message-out-handler : socket close"
			continue
		}
		cmd := <-c.messagesQueue
		c.socket.Write([]byte(cmd))
	}
}

func (c *Communicator) errorListener() {
	for {
		msg := <-c.errors
		fmt.Println(msg)
	}
}

func (c *Communicator) SendMessage(cmdName string, params ...interface{}) {
	c.messagesQueue <- wrapMessage(ircCommands[cmdName].format, cmdName, params...)
}

func (c *Communicator) Run() {
	go c.handleSendMessage()
	go c.errorListener()
}
