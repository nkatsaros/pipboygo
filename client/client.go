package client

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/nkatsaros/pipboygo/protocol"
)

// ErrBusy lets the user know the server is busy
var ErrBusy = errors.New("the server is busy")

const port = 27000

var pong = []byte{0, 0, 0, 0, 0}

type Command struct {
	Type int           `json:"type"`
	Args []interface{} `json:"args"`
	ID   int           `json:"id"`
}

type Event struct {
	Channel int
	Data    interface{}
}

type Client struct {
	conn net.Conn

	logger *log.Entry

	event     chan *Event
	err       chan error
	busyCheck chan error

	wg   *sync.WaitGroup
	done chan struct{}
}

func New(ip string) (c *Client, err error) {
	c = &Client{
		event:     make(chan *Event),
		err:       make(chan error),
		busyCheck: make(chan error),

		wg:   &sync.WaitGroup{},
		done: make(chan struct{}),
	}

	c.logger = log.WithFields(log.Fields{
		"component": "client",
		"server_ip": ip,
	})

	c.logger.Info("connecting")

	c.conn, err = net.Dial("tcp4", net.JoinHostPort(ip, strconv.Itoa(port)))
	if err != nil {
		return nil, err
	}

	c.logger.Info("connected")

	c.wg.Add(1)
	go c.main()

	// check if the first message we receive yields a busy or error
	select {
	case err = <-c.busyCheck:
	case err = <-c.err:
	}
	if err != nil {
		c.logger.Error(err)
		c.Close()
		return nil, err
	}

	return c, nil
}

func (c *Client) main() {
	defer c.wg.Done()

	firstMessage := true

	decoder := protocol.NewDecoder(c.conn)
	for {
		select {
		case <-c.done:
			return
		default:
		}

		channel, data, err := decoder.Decode()
		if cErr, ok := err.(*protocol.UnknownChannelError); ok {
			c.logger.WithFields(log.Fields{
				"channel": cErr.Channel,
				"data":    fmt.Sprintf("%X", cErr.Data),
			}).Error(cErr)
			select {
			case <-c.done:
				return
			case c.err <- err:
			}
			continue
		} else if err != nil {
			select {
			case <-c.done:
				return
			case c.err <- err:
			}
			continue
		}

		// check if the first message received is a busy message
		if firstMessage {
			firstMessage = false
			if channel == protocol.ChannelBusy {
				select {
				case <-c.done:
					return
				case c.busyCheck <- ErrBusy:
				}
			} else {
				close(c.busyCheck)
			}
		}

		switch channel {
		case protocol.ChannelHeartbeat:
			c.conn.Write(pong)
		default:
			select {
			case <-c.done:
				return
			case c.event <- &Event{channel, data}:
			}
		}
	}
}

func (c *Client) Err() <-chan error {
	return c.err
}

func (c *Client) Event() <-chan *Event {
	return c.event
}

func (c *Client) Close() (err error) {
	close(c.done)
	err = c.conn.Close()

	c.wg.Wait()

	close(c.event)
	close(c.err)

	return err
}
