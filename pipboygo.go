package main

import (
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nkatsaros/pipboygo/autodiscovery"
	"github.com/nkatsaros/pipboygo/client"
)

func main() {
	logger := log.WithFields(log.Fields{
		"component": "app",
	})

	// watch for Ctrl+C
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

loop:
	for {
		// search for a non-busy game
		a, err := autodiscovery.Listen()
		if err != nil {
			logger.Fatalln(err)
		}
		var game autodiscovery.Game
	autodiscoverLoop:
		for {
			select {
			case game = <-a.Games():
				if game.IsBusy == false {
					break autodiscoverLoop
				}
			case <-signals:
				break loop
			}
		}
		a.Close()

		// connect
		c, err := client.New(game.IP)
		if err != nil {
			// bad or busy server, try again in 1 second
			logger.Error(err)
			select {
			case <-time.After(1 * time.Second):
			case <-signals:
				break loop
			}
			continue
		}

		// do nothing with the messages we receive!
		err = func() error {
			defer c.Close()
			for {
				select {
				case event := <-c.Event():
					logger.WithField("channel", event.Channel).Info("received")
				case err := <-c.Err():
					logger.Error(err)
					return err
				case <-signals:
					return nil
				}
			}
		}()
		if err != nil {
			logger.Error(err)
			continue
		}

		break
	}
}
