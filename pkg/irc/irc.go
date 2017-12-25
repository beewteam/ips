package irc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type EventCallback func(eventName string)

type ResponseCallback func(response string, err string)
type Communicator struct {
	socket  net.Conn
	pMsg    message
	mQueue  []message
	msgs    chan message
	workers []*Worker

	log             *logrus.Logger
	eventDispatcher map[string]eventCallbackList
	OnMsg           func(response string)
	OnError         func(err string)
}

type eventCallbackList []EventCallback
type ircCommand struct {
	format string
}
type message struct {
	data string
	done bool
}

const (
	tCmdName       = "@@CmdName@@"
	protoDelim     = "\r\n"
	protoMaxLength = 512

	warningError = "1"

	exitCtrl = 1
	timeout  = 10 * time.Second
)

var (
	ircCommands map[string]ircCommand
)

type Worker struct {
	in  chan []byte
	out chan []byte
	err chan error
	ctl chan bool
}

func NewWorker(in chan []byte, err chan error, ctl chan bool) *Worker {
	return &Worker{
		in:  in,
		out: make(chan []byte),
		err: err,
		ctl: ctl,
	}
}

func (c *Communicator) Subscribe(eventName string, callbackFunction EventCallback) (err error) {
	if callList, ok := c.eventDispatcher[eventName]; ok {
		c.eventDispatcher[eventName] = append(callList, callbackFunction)
		err = errors.New("Communicator.Subscribe: cannot add callback to not existing event")
	}
	return err
}

func (c *Communicator) Init() {
	ircCommands = map[string]ircCommand{
		"USER":    {tCmdName + " %s %s %s :%s"},
		"PING":    {tCmdName + " %s"},
		"PART":    {tCmdName + " %s"},
		"NICK":    {tCmdName + " %s"},
		"QUIT":    {tCmdName + " :%s"},
		"TOPIC":   {tCmdName + " %s :%s"},
		"PRIVMSG": {tCmdName + " %s :%s"},
	}
	c.msgs = make(chan message)
	c.mQueue = make([]message, 0)

	c.eventDispatcher = make(map[string]eventCallbackList)
	c.eventDispatcher["*"] = make([]EventCallback, 0)

	c.log = logrus.New()
	c.log.Out = ioutil.Discard

}

func (c *Communicator) SendMessage(cmdName string, params ...interface{}) {
	go func() {
		c.msgs <- message{
			wrapMessage(ircCommands[cmdName].format, cmdName, params...),
			false,
		}
	}()
}

func (c *Communicator) SetLog(logPath string) error {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err = os.Mkdir("log", 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filepath.Join("log", logPath), os.O_CREATE|os.O_WRONLY, 0655)
	if err == nil {
		c.log.Out = file

	}
	return err
}

func (c *Communicator) Close() {
	for _, v := range c.workers {
		v.ctl <- true
	}
	c.socket.Close()
}

func wrapMessage(format string, cmdName string, params ...interface{}) string {
	parsedFormat := strings.Replace(format, tCmdName, cmdName, 1) + protoDelim
	return fmt.Sprintf(parsedFormat, params...)
}

func (c *Communicator) Run(hostname string, port string) (err error) {
	c.socket, err = net.Dial("tcp", hostname+":"+port)
	if err == nil {
		writerChan := make(chan []byte)
		errorsChan := make(chan error)
		ctlR := make(chan bool, 1)
		ctlW := make(chan bool, 1)

		var readerW = NewWorker(nil, errorsChan, ctlR)
		var writerW = NewWorker(writerChan, errorsChan, ctlW)

		c.workers = append(c.workers, readerW)
		c.workers = append(c.workers, readerW)

		go reader(c, readerW)
		go writer(c, writerW)
	}

	return err
}

func reader(c *Communicator, w *Worker) {
	var buf = make([]byte, protoMaxLength)
	var rerr error

	for {
		select {
		case <-w.ctl:
			return
		case w.err <- rerr:
			if rerr != nil {
				c.log.Info(rerr.Error())
				close(w.err)
				return
			}
		default:
			c.socket.SetReadDeadline(time.Now().Add(timeout))
			_, rerr = c.socket.Read(buf)
			if rerr != nil {
				return
			}

			reply := string(buf)
			reply = strings.TrimSuffix(reply, "\r")
			reply = strings.TrimSuffix(reply, "\n")

			for _, v := range c.eventDispatcher {
				// Check if event happens
				if false {
					for _, f := range v {
						f(reply)
					}
				}
			}
			for _, f := range c.eventDispatcher["*"] {
				f(reply)
			}
		}
	}
}

func writer(c *Communicator, w *Worker) {
	var werr error

	for {
		select {
		case <-w.ctl:
			return
		case w.err <- werr:
			if werr != nil {
				c.log.Info(werr.Error())
				close(w.err)
				return
			}
		case data := <-w.in:
			c.socket.SetWriteDeadline(time.Now().Add(timeout))
			_, werr = c.socket.Write(data)
		}
	}
}
