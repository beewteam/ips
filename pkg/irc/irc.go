package irc

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type (
	Communicator struct {
		socket    net.Conn
		pMsg      message
		mQueue    []message
		msgs      chan message
		readerIn  chan []byte
		writerOut chan []byte
		errors    chan error
		control   chan int

		OnMsg   func(response string)
		OnError func(err string)
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
		done bool
	}
)

const (
	tCmdName       = "@@CmdName@@"
	protoDelim     = "\r\n"
	protoMaxLength = 512

	warningError = "1"

	exitCtrl = 1
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
	c.control = make(chan int)
	c.msgs = make(chan message)

	c.mQueue = make([]message, 0)

	c.pMsg.done = true
}

func (c *Communicator) SendMessage(cmdName string, callback ResponseCallback, params ...interface{}) {
	go func() {
		c.msgs <- message{
			wrapMessage(ircCommands[cmdName].format, cmdName, params...),
			callback,
			false,
		}
	}()
}

func (c *Communicator) Close() {
	c.control <- exitCtrl
	c.control <- exitCtrl
	c.control <- exitCtrl

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
		go router(c, c.control, c.msgs, c.writerOut, c.readerIn, c.errors)
	}

	return err
}

func (c *Communicator) send(writer chan<- []byte) {
	if c.pMsg.done && len(c.mQueue) != 0 {
		c.pMsg = c.mQueue[0]
		c.mQueue = c.mQueue[1:]
		go func() {
			writer <- []byte(c.pMsg.data)
		}()
	}
}

func router(c *Communicator, control <-chan int, messages <-chan message, writer chan<- []byte, reader <-chan []byte, err <-chan error) {
	var buf = make([]byte, protoMaxLength)

	for {
		select {
		case ctl := <-control:
			fmt.Println("router ctl: " + strconv.FormatInt(int64(ctl), 10))
			if ctl == exitCtrl {
				return
			}
		case buf = <-reader:
			reply := string(buf)
			reply = strings.TrimSuffix(reply, "\n")
			reply = strings.TrimSuffix(reply, "\r")

			c.OnMsg(reply)
			c.pMsg.done = true
			c.send(writer)
		case e := <-err:
			if e != nil {
				fmt.Println(e)
			}
		case msg := <-messages:
			c.mQueue = append(c.mQueue, msg)
			c.send(writer)
		}
	}
}

func reader(socket net.Conn, control <-chan int, out chan<- []byte, err chan<- error) {
	var buf []byte
	var rerr error

	for {
		buf = make([]byte, protoMaxLength)
		_, rerr = socket.Read(buf)
		select {
		case ctl := <-control:
			fmt.Println("writer ctl: " + strconv.FormatInt(int64(ctl), 10))
			if ctl == exitCtrl {
				return
			}
		case out <- buf:
		case err <- rerr:
		}
	}
}

func writer(socket net.Conn, control <-chan int, in <-chan []byte, err chan<- error) {
	var werr error

	for {
		select {
		case ctl := <-control:
			fmt.Println("reader ctl: " + strconv.FormatInt(int64(ctl), 10))
			if ctl == exitCtrl {
				return
			}
		case err <- werr:
		case data := <-in:
			_, werr = socket.Write(data)

		}
	}
}
