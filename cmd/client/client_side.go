package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/beewteam/ips/cmd/client/ui"
	"github.com/beewteam/ips/pkg/irc"
	tui "github.com/marcusolsson/tui-go"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

type Server struct {
	Hostname string
	Port     string
}

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

	com irc.Communicator

	ui   tui.UI
	desk map[string]tui.Widget

	rl *readline.Instance
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
	// if len(params) != 2 {
	// 	output = append(output, "Error: Wrong numbers of arguments")
	// 	return
	// }

	// c.Server = Server{
	// 	Hostname: params[0],
	// 	Port:     params[1],
	// }
	// if c.Server.Connect() {
	// 	output = append(output, "Connected to "+c.Server.String())
	// } else {
	// 	output = append(output, "Error: Cannot connect to "+c.Server.String())
	// }

	// return
	return []string{""}
}

// Syntax: help
func help(c *Client, params []string) (output []string) {
	output = append(output, "Help menu")
	for _, cmd := range clientCommands {
		output = append(output, cmd.name+"["+cmd.shortcut+"]"+": "+cmd.desc)
	}
	return
}

func (c *Client) setupUI() {
	channelList := tui.NewList()
	sidebar := tui.NewVBox(
		tui.NewLabel("Channels:"),
		channelList,
		tui.NewSpacer())
	sidebar.SetBorder(true)
	c.desk["channel_list"] = channelList

	// ChatArea
	chat := ui.NewChatArea()

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.OnSubmit(func(entry *tui.Entry) {
		chat.AddNewMessage(entry.Text())
		entry.SetText("")
	})

	inputBox := tui.NewHBox(
		tui.NewLabel("user>"),
		input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	mainView := tui.NewVBox(
		chat,
		inputBox)

	statusbar := tui.NewStatusBar("")
	statusbarBox := tui.NewHBox(
		tui.NewLabel("Statusbar:"),
		statusbar)
	statusbarBox.SetBorder(true)
	statusbarBox.SetSizePolicy(tui.Expanding, tui.Maximum)
	c.desk["statusbar"] = statusbar

	root := tui.NewVBox(
		tui.NewHBox(
			sidebar,
			mainView),
		statusbar)

	c.ui = tui.New(root)
	c.ui.SetKeybinding("Esc", func() { c.ui.Quit() })
	//c.ui.SetKeybinding("Up", func() { msgArea.Scroll(0, 1) })
	//c.ui.SetKeybinding("Down", func() { msgArea.Scroll(0, -1) })
}

func NewClient(configFile string) *Client {
	var settings UserConfig
	rl, err := readline.New("irc> ")
	if err != nil {
		fmt.Println(color.RedString("client") + ": cannot init input interface")
		panic(err)
	}
	defer rl.Close()

	// fmt.Printf("Trying to load %s\n", color.GreenString(configFile))
	// if _, err = os.Stat(configFile); err != nil {
	// 	settings, err = ParseConfigFile(configFile)
	// 	if err != nil {
	// 		fmt.Println(color.RedString("client") + ": cannot load config")
	// 		panic(err)
	// 	}
	// }

	_ = settings

	// var client = irc.Client{
	// 	Account: irc.Account{
	// 		Username: settings.UserData.Username,
	// 		FullName: settings.UserData.FullName,
	// 	},
	// 	Server: irc.Server{
	// 		Hostname: settings.ServerData.Hostname,
	// 		Port:     settings.ServerData.Port},
	// }

	c := &Client{
		desk: make(map[string]tui.Widget),
	}

	c.setupUI()

	return c
}

func (c *Client) Run() error {
	return c.ui.Run()
	//client.Run()
	// var com irc.Communicator
	// com.Init()
	// defer com.Close()
	// errlog := com.SetLog("irc.log")
	// if errlog != nil {
	// 	fmt.Println("Cannot setup log file")
	// 	return
	// }
	// com.Subscribe("*", printM)

	// err := com.Run("irc.freenode.com", "8000")
	// if err != nil {
	// 	fmt.Printf("Cannot run communicator\n")
	// 	return
	// }

	// for {
	// 	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	// 	if text[0] == 'q' {
	// 		return
	// 	} else if text[0] == 'p' {
	// 		com.SendMessage("PING", print, "irc.freenode.com")
	// 	} else if text[0] == 'u' {
	// 		com.SendMessage("USER", "testion", "0", "*", "Oliak")
	// 	}
	// }
	return nil
}

func init() {
	clientCommands = []command{
		{"help", "h", "print help menu", help},
		{"connect", "c", "connect client to server", connect},
	}
}

/*func Reg(c net.Conn, nick string, username string, fullname string) bool {
	fmt.Fprintf(
		c,
		"USER %s 0 0 :%s\r\n",
		nick, username, fullname)
	return true
}*/
