package main

import (
	"encoding/binary"
	"encoding/json"
	"image/png"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nkatsaros/pipboygo/autodiscovery"
	"github.com/nkatsaros/pipboygo/protocol"
)

const port = 27000

var pong = []byte{0, 0, 0, 0, 0}

type Command struct {
	Type int           `json:"type"`
	Args []interface{} `json:"args"`
	ID   int           `json:"id"`
}

func connect(logger *log.Entry, game autodiscovery.Game) (err error) {
	conn, err := net.Dial("tcp4", net.JoinHostPort(game.IP, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		time.Sleep(1 * time.Second)
		data, _ := json.Marshal(Command{13, []interface{}{}, 1})
		size := uint32(len(data))
		binary.Write(conn, binary.LittleEndian, size)
		conn.Write([]byte{0x5})
		conn.Write(data)
	}()

	decoder := protocol.NewDecoder(conn)
	for {
		channel, data, err := decoder.Decode()
		if err != nil {
			return err
		}

		switch channel {
		case 0:
			conn.Write(pong)
		case 4:
			lm := data.(protocol.LocalMap)
			f, err := os.OpenFile("image.png", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			err = png.Encode(f, lm.Image)
			if err != nil {
				return err
			}
			f.Close()
		default:
			logger.WithField("channel", channel).Info(data)
		}
	}

	return nil
}

func main() {
	for {
		a, err := autodiscovery.Listen()
		if err != nil {
			log.Fatalln(err)
		}
		var game autodiscovery.Game
		for game = range a.Games() {
			if game.IsBusy == false {
				break
			}
		}
		a.Close()

		logger := log.WithFields(log.Fields{
			"game_ip":   game.IP,
			"game_type": game.MachineType,
		})
		logger.Info("connecting")
		if err = connect(logger, game); err != nil {
			continue
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
}
