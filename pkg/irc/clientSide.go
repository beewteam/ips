package irc

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type Account struct {
	Username string
	FullName string
	Nick     string
}

type Client struct {
	Account  Account
	Channel  string
	Data     string
	Server   Server
	isActive bool
}

type command struct {
	name     string
	shortcut string
	desc     string
	handler  func(c *Client, params []string) []string
}

var (
	clientCommands []command
)

// func (c *Client) Login(nick string) bool {
// 	return Reg(c.Server.Conn, nick, c.Account.Username, c.Account.FullName)
// }

// func (c *Client) JoinChannel(channel string) bool {
// 	c.Channel = channel
// 	return Join(c.Server.Conn, channel)
// }

// func (c *Client) LeaveChannel() bool {
// 	return Part(c.Server.Conn, c.Channel)
// }

// func (c *Client) HandleData() bool {
// 	message, _ := bufio.NewReader(c.Server.Conn).ReadString('\r')
// 	if len(message) > 0 {
// 		fmt.Printf("Message:%s\n", message)
// 	}
// 	return true
// }

// func (c *Client) LogMessage(nick string, msg string) bool {
// 	fmt.Printf(
// 		"%s - [%s] <%s> %s\n",
// 		time.Now().String(), c.Channel, c.Account.Nick, msg)
// 	return true
// }

// func (c *Client) Close() bool {
// 	c.Server.Conn.Close()
// 	return true
// }

// func (c *Client) Auth(passwd string) bool {
// 	return Msg(c.Server.Conn, "NickServ", "IDENTIFY "+passwd)
// }

func findHandler(commandName string) *command {
	for i := range clientCommands {
		if clientCommands[i].name == commandName || clientCommands[i].shortcut == commandName {
			return &clientCommands[i]
		}
	}
	return &clientCommands[1]
}

func (c *Client) handleCommand(cmdLine string) (out []string) {
	word := strings.Fields(cmdLine)
	cmd := findHandler(strings.ToLower(word[0]))
	return cmd.handler(c, word[1:])
}

func (c *Client) addListener(wg *sync.WaitGroup, reader *bufio.Reader, writer *bufio.Writer) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for c.isActive {
			var out []string

			cmdLine, err := reader.ReadString('\n')
			if err != nil {
				// Handle this by logging to file
				continue
			}

			out = c.handleCommand(cmdLine)

			_, err = writer.WriteString(strings.Join(out, "\n") + "\n")
			if err != nil {
				// Handle this by logging to file
				continue
			}
			writer.Flush()
		}
	}()
}

func AddUserCommandListener(wg sync.WaitGroup, reader io.Reader) {

}

// Syntax: connect server-name port-number
func connect(c *Client, params []string) (output []string) {
	if len(params) != 2 {
		output = append(output, "Error: Wrong numbers of arguments")
		return
	}

	c.Server = Server{
		Hostname: params[0],
		Port:     params[1],
	}
	if c.Server.Connect() {
		output = append(output, "Connected to "+c.Server.String())
	} else {
		output = append(output, "Error: Cannot connect to "+c.Server.String())
	}

	return
}

// Syntax: help
func help(c *Client, params []string) (output []string) {
	output = append(output, "Help menu")
	for _, cmd := range clientCommands {
		output = append(output, cmd.name+"["+cmd.shortcut+"]"+": "+cmd.desc)
	}
	return
}

func ListAvaiableServers(c *Client, msg string) {

}

func (c *Client) Run() bool {
	var wg sync.WaitGroup

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	c.addListener(&wg, r, w)

	wg.Wait()
	return true
}

func (c *Client) Init() bool {
	clientCommands = []command{
		{"help", "h", "print help menu", help},
		{"connect", "c", "connect client to server", connect},
	}
	c.isActive = true
	return true
}

/*func Reg(c net.Conn, nick string, username string, fullname string) bool {
	fmt.Fprintf(
		c,
		"USER %s 0 0 :%s\r\n",
		nick, username, fullname)
	return true
}*/
