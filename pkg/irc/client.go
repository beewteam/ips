package irc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Account struct {
	Username string
	FullName string
	Nick     string
}

type Client struct {
	Account Account
	Channel string
	Data    string
	Server  Server
}

type command struct {
	name     string
	shortcut string
	handler  func(c *Client, params []string) (bool, []string)
}

var (
	clientCommands []command
)

func (c *Client) Login(nick string) bool {
	return Reg(c.Server.Conn, nick, c.Account.Username, c.Account.FullName)
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
		time.Now().String(), c.Channel, c.Account.Nick, msg)
	return true
}

func (c *Client) Close() bool {
	c.Server.Conn.Close()
	return true
}

func (c *Client) Auth(passwd string) bool {
	return Msg(c.Server.Conn, "NickServ", "IDENTIFY "+passwd)
}

func findHandler(commandName string) *command {
	for i := range clientCommands {
		if clientCommands[i].name == commandName || clientCommands[i].shortcut == commandName {
			return &clientCommands[i]
		}
	}
	return &clientCommands[1]
}

func (c *Client) handleCommand(reader *bufio.Reader, writer *bufio.Writer) (bool, bool) {
	var s bool
	var out []string

	cmdLine, err := reader.ReadString('\n')
	if err == nil {
		word := strings.Fields(cmdLine)
		s, out = findHandler(word[0]).handler(c, word[1:])
		for i := range out {
			fmt.Println(out[i])
			//_, err = writer.WriteString(out[i] + "\n")
			//if err != nil {
			//}
		}
	} else {

	}

	return false, s
}

func (c *Client) AddListener(wg *sync.WaitGroup, reader *bufio.Reader, writer *bufio.Writer) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		isActive := true
		for isActive {
			err, stop := c.handleCommand(reader, writer)
			if err {

			}
			if stop {
				isActive = false
			}
		}
	}()
}

func AddUserCommandListener(wg sync.WaitGroup, reader io.Reader) {

}

// Syntax: connect server-name port-number
func connect(c *Client, params []string) (successful bool, output []string) {
	if len(params) != 2 {
		output = append(output, "Error: Wrong numbers of arguments")
		successful = false
		return
	}

	c.Server = Server{
		Hostname: params[0],
		Port:     params[1],
	}
	if c.Server.Connect() {
		output = append(output, "Connected to "+c.Server.String())
		successful = true
	} else {
		output = append(output, "Error: Cannot connect to "+c.Server.String())
		successful = false
	}

	return
}

// Syntax: connect server-name port-number
func help(c *Client, params []string) (successful bool, output []string) {
	output = append(output, "help")
	successful = true
	return
}

func ListAvaiableServers(c *Client, msg string) {

}

func (c *Client) Run() bool {
	var wg sync.WaitGroup

	//if !client.Login(settings.UserData.Nickname) {
	//	os.Exit(1)
	//}

	//client.Auth(string(settings.UserData.Password))

	//client.JoinChannel(settings.UserData.Chat)

	//for client.HandleData() {
	//}

	//client.Close()

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	c.AddListener(&wg, r, w)

	wg.Wait()
	return true
}

func (c *Client) Init() bool {
	clientCommands = []command{
		{"HELP", "H", help},
		{"CONNECT", "C", connect},
	}
	return true
}
