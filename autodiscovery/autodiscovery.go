package autodiscovery

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

const port = 28000
const frequency = 1 * time.Second
const timeout = 1 * time.Second
const payload = `{"cmd": "autodiscover"}`

var listenAddr = &net.UDPAddr{IP: net.IPv4zero, Port: port}
var broadcastAddr = &net.UDPAddr{IP: net.IPv4bcast, Port: port}

type Game struct {
	IP          string
	IsBusy      bool
	MachineType string
}

type Autodiscovery struct {
	conn  *net.UDPConn
	games chan Game

	wg   *sync.WaitGroup
	done chan struct{}
}

func Listen() (a *Autodiscovery, err error) {
	a = &Autodiscovery{
		games: make(chan Game),
		wg:    &sync.WaitGroup{},
		done:  make(chan struct{}),
	}

	a.conn, err = net.ListenUDP("udp4", listenAddr)
	if err != nil {
		return nil, err
	}

	a.wg.Add(2)
	go a.discover()
	go a.broadcastAutodiscover()

	return a, nil
}

func (a *Autodiscovery) Games() <-chan Game {
	return a.games
}

func (a *Autodiscovery) Close() error {
	close(a.done)
	a.wg.Wait()
	close(a.games)
	return a.conn.Close()
}

func (a *Autodiscovery) discover() {
	defer a.wg.Done()

	data := make([]byte, 1472)
	for {
		select {
		case <-a.done:
			return
		default:
		}

		var game Game
		a.conn.SetReadDeadline(time.Now().Add(timeout))
		n, addr, err := a.conn.ReadFromUDP(data)
		if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
			continue
		} else if err != nil {
			log.Error(err)
			continue
		}
		if err = json.Unmarshal(data[:n], &game); err != nil {
			log.Error(err)
			continue
		}
		if game.MachineType == "" {
			continue
		}

		game.IP = addr.IP.String()
		a.games <- game
	}
}

func (a *Autodiscovery) broadcastAutodiscover() {
	defer a.wg.Done()

	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-a.done:
			return
		case <-timer.C:
			_, err := a.conn.WriteToUDP([]byte(payload), broadcastAddr)
			if err != nil {
				log.Error(err)
			}
			timer.Reset(frequency)
		}
	}
}
