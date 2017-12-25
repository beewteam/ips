package main

import (
	"strings"
)

type command struct {
	name     string
	shortcut string
	desc     string
	handler  func(c *Client, param []string) bool
}

var (
	clientCommands []command
)

func incorrectCmd(c *Client, cmd string) {
	c.chat.AddNewMessage("Incorrect command:" + cmd)
}

// Syntax: help
func help(c *Client, param []string) bool {
	output := make([]string, 0)
	output = append(output, "Help menu")
	for _, cmd := range clientCommands {
		output = append(output, cmd.name+"["+cmd.shortcut+"]"+": "+cmd.desc)
	}
	c.chat.AddNewMessage(strings.Join(output, "\n"))
	return true
}

func post(c *Client, param []string) bool {
	if !c.IsRegistered {
		return false
	}
	if len(param) != 2 {
		return false
	}
	c.com.SendMessage("PRIVMSG", param[0], param[1])
	return true
}

func join(c *Client, param []string) bool {
	if len(param) != 2 {
		return false
	}
	if !c.IsRegistered {
		return false
	}

	c.channelBar.AddChannel(param[0])
	c.com.SendMessage("JOIN", param[0])
	return true
}

func changeNick(c *Client, param []string) bool {
	if len(param) != 2 {
		return false
	}

	c.com.SendMessage("NICK", param[0])
	return true
}

func regUser(c *Client, param []string) bool {
	c.com.SendMessage("NICK", c.Account.Nick)
	c.com.SendMessage("USER", c.Account.Username, "8", "*", c.Account.FullName)
	c.IsRegistered = true
	return true
}

func serverInfo(c *Client, param []string) bool {
	c.chat.AddNewMessage("Server hostname:" + c.Server.Hostname + ", port:" + c.Server.Port)
	return true
}

func init() {
	clientCommands = []command{
		{"help", "h", "print help menu", help},
		{"join", "j", "join to chat", join},
		{"post", "p", "post message in chat", post},
		{"change-nick", "cn", "changes nick", changeNick},
		{"register", "r", "register user", regUser},
		{"server-info", "si", "print server info", serverInfo},
	}
}
