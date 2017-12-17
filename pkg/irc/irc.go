package irc

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type (
	Communicator struct {
		socket        net.Conn
		processingMsg message
		msgs          chan message
		readerIn      chan []byte
		writerOut     chan []byte
		errors        chan error
		control       chan int
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
}

func (c *Communicator) SendMessage(cmdName string, callback ResponseCallback, params ...interface{}) {
	go func() {
		c.msgs <- message{
			wrapMessage(ircCommands[cmdName].format, cmdName, params...),
			callback,
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
		go router(c.control, c.msgs, c.writerOut, c.readerIn, c.errors)
	}

	return err
}

func router(control <-chan int, messages <-chan message, writer chan<- []byte, reader <-chan []byte, err <-chan error) {
	var status int
	var msg message
	var buf = make([]byte, protoMaxLength)

	for {
		select {
		case ctl := <-control:
			if ctl == 1 {
				fmt.Println("router ctl: " + strconv.FormatInt(int64(ctl), 10))
				return
			}
		case e := <-err:
			fmt.Println(e)
			return
		default:
			if status == 0 {
				select {
				case msg = <-messages:
					status = 1
				default:
				}
			} else if status == 1 {
				select {
				case writer <- []byte(msg.data):
					status = 0
				default:
				}
			}

			select {
			case buf = <-reader:
				fmt.Println(string(buf))
			default:
			}
			//fmt.Println("idle")
		}
	}
}

func reader(socket net.Conn, control <-chan int, out chan<- []byte, err chan<- error) {
	buf := make([]byte, protoMaxLength)
	var state int
	var rerr error

	for {
		select {
		case ctl := <-control:
			fmt.Println("writer ctl: " + strconv.FormatInt(int64(ctl), 10))
			if ctl == exitCtrl {
				return
			}
		default:
		}

		if state == 0 {
			_, rerr = socket.Read(buf)
			state = 1
			// For debug
			fmt.Println(string(buf))
		} else if state == 1 {
			select {
			case err <- rerr:
				rerr = nil
				return
			case out <- buf:
				state = 1
			default:
			}
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
		default:
		}

		if werr == nil {
			select {
			case data := <-in:
				_, werr = socket.Write(data)
			default:
			}
		} else {
			select {
			case err <- werr:
				werr = nil
			default:
			}
		}
	}
}
