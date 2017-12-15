package irc

import (
	"fmt"
	"io"
	"net"
	"strings"
)

type (
	Communicator struct {
		socket        net.Conn
		processingMsg message
		readerIn      chan []byte
		writerOut     chan []byte
		errors        chan error
		control       chan byte
	}

	ResponseCallback func(response string, err string)
)

type (
	ircCommand struct {
		format string
	}

	message struct {
		data string
		f    ResponseCallback
	}
)

const (
	tCmdName       = "@@CmdName@@"
	protoDelim     = "\r\n"
	protoMaxLength = 512

	warningError = "1"
)

var (
	ircCommands map[string]ircCommand
)

func (c *Communicator) Init() {
	ircCommands = map[string]ircCommand{
		"PING":    {tCmdName + " %s"},
		"PART":    {tCmdName + " %s"},
		"NICK":    {tCmdName + " %s"},
		"QUIT":    {tCmdName + " :%s"},
		"TOPIC":   {tCmdName + " %s :%s"},
		"PRIVMSG": {tCmdName + " %s :%s"},
	}

	c.readerIn = make(chan []byte)
	c.writerOut = make(chan []byte)
	c.errors = make(chan error)
	c.control = make(chan byte)
}

func (c *Communicator) SendMessage(cmdName string, f ResponseCallback, params ...interface{}) {
	//go func() {
	//c.messagesQueue <- message{
	//	wrapMessage(ircCommands[cmdName].format, cmdName, params...),
	//	f,
	//}
	//}()
}

func (c *Communicator) Close() {
	c.socket.Close()
}

func wrapMessage(format string, cmdName string, params ...interface{}) string {
	parsedFormat := strings.Replace(format, tCmdName, cmdName, 1) + protoDelim
	return fmt.Sprintf(parsedFormat, params...)
}

func (c *Communicator) Run(hostname string, port string) (err error) {
	c.socket, err = net.Dial("tcp", hostname+":"+port)
	if err == nil {
		go reader(c.socket, c.control, c.readerIn, c.errors)
		go writer(c.socket, c.control, c.writerOut, c.errors)
		go router(c.errors)
	}

	return err
}

func router(err <-chan error) {
	for {
		select {
		case e := <-err:
			fmt.Println(e)
		default:
			fmt.Println("idle")
		}
	}
}

func reader(socket net.Conn, control <-chan byte, out chan<- []byte, err chan<- error) {
	buf := make([]byte, protoMaxLength)
	for {
		select {
		case ctl := <-control:
			fmt.Println("writer ctl: " + string(ctl))
		default:
			_, rerr := socket.Read(buf)
			if rerr == nil && rerr != io.EOF {
				out <- buf
			}
			err <- rerr
		}
	}

}

func writer(socket net.Conn, control <-chan byte, in <-chan []byte, err chan<- error) {
	for {
		select {
		case ctl := <-control:
			fmt.Println("reader ctl: " + string(ctl))
		default:
			_, werr := socket.Write(<-in)
			err <- werr
		}
	}
}
