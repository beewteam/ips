package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/beewteam/ips/cmd/client/ui"
	"github.com/beewteam/ips/pkg/irc"
	tui "github.com/marcusolsson/tui-go"

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
	Account      Account
	Channel      string
	Data         string
	Server       Server
	IsRegistered bool

	com irc.Communicator

	ui         tui.UI
	chat       *ui.ChatArea
	channelBar *ui.ChannelBar
}

func findHandler(commandName string) *command {
	for i := range clientCommands {
		if clientCommands[i].name == commandName || clientCommands[i].shortcut == commandName {
			return &clientCommands[i]
		}
	}
	return &clientCommands[1]
}

func (c *Client) handleUserInput(line string) {
	words := strings.Fields(line)
	cmd := findHandler(strings.ToLower(words[0]))
	cmd.handler(c, words[1:])
}

func (c *Client) setupUI() {
	// Channel sidebar
	c.channelBar = ui.NewChannelBar()

	// ChatArea
	c.chat = ui.NewChatArea()

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	input.OnSubmit(func(entry *tui.Entry) {
		userInputLine := entry.Text()
		if len(userInputLine) > 0 {
			c.handleUserInput(userInputLine)
			c.chat.AddNewMessage(userInputLine)
		}
		entry.SetText("")
	})

	inputBox := tui.NewHBox(
		tui.NewLabel("user>"),
		input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	mainView := tui.NewVBox(
		c.chat.ToWidget(),
		inputBox)

	statusbar := tui.NewStatusBar("")
	statusbarBox := tui.NewHBox(
		tui.NewLabel("Statusbar:"),
		statusbar)
	statusbarBox.SetBorder(true)
	statusbarBox.SetSizePolicy(tui.Expanding, tui.Minimum)

	root := tui.NewVBox(
		tui.NewHBox(
			c.channelBar.ToWidget(),
			mainView),
		statusbar)

	c.ui = tui.New(root)
	c.ui.SetKeybinding("Esc", func() { c.ui.Quit() })
	//c.ui.SetKeybinding("Up", func() { msgArea.Scroll(0, 1) })
	//c.ui.SetKeybinding("Down", func() { msgArea.Scroll(0, -1) })
}

func NewClient(configFile string) (*Client, error) {
	var settings UserConfig

	fmt.Printf("Trying to load %s\n", color.GreenString(configFile))
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		settings, err = ParseConfigFile(configFile)
		if err != nil {
			fmt.Println(color.RedString("client") + ": cannot load config")
			return nil, err
		}
	} else {
		return nil, err
	}

	c := &Client{
		Account: Account{
			Username: settings.UserData.Username,
			FullName: settings.UserData.FullName,
		},
		Server: Server{
			Hostname: settings.ServerData.Hostname,
			Port:     settings.ServerData.Port,
		},
		IsRegistered: false,
	}
	c.setupUI()

	return c, nil
}

func (c *Client) postMessageInChat(msg string) {
	c.chat.AddNewMessage(msg)
}

func (c *Client) Run() error {
	c.com.Init()

	errlog := c.com.SetLog("irc.log")
	if errlog != nil {
		fmt.Println("Cannot setup log file")
		return errlog
	}
	c.com.Subscribe("*", c.postMessageInChat)

	err := c.com.Run(c.Server.Hostname, c.Server.Port)
	if err != nil {
		fmt.Println("Cannot run communicator")
		return err
	}

	return c.ui.Run()
}

func (c *Client) Close() {
	c.com.Close()
}
