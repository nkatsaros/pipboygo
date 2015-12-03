package main

import (
	"flag"
	"image"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nkatsaros/pipboygo/autodiscovery"
	"github.com/nkatsaros/pipboygo/client"
	"github.com/nkatsaros/pipboygo/protocol"
)

var publicFlag = flag.String("public", "", "path to public files")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type memresponse struct {
	Type   string                `json:"type"`
	Memory protocol.PipboyMemory `json:"memory"`
}

type localmapmetaresponse struct {
	Type    string           `json:"type"`
	Extents protocol.Extents `json:"extents"`
}

func main() {
	flag.Parse()

	logger := log.WithFields(log.Fields{
		"component": "app",
	})

	_ = logger
	if *publicFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error(err)
			return
		}

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
			}
		}
		a.Close()

		// connect
		pc, err := client.New(game.IP)
		if err != nil {
			logger.Error(err)
			return
		}

		err = func() error {
			defer pc.Close()

			localmaptimer := time.NewTimer(0)
			defer localmaptimer.Stop()
			for {
				select {
				case event := <-pc.Event():
					switch t := event.Data.(type) {
					case protocol.LocalMap:
						conn.WriteJSON(localmapmetaresponse{
							Type:    "local_map_metadata",
							Extents: t.Extents,
						})
						conn.WriteMessage(websocket.BinaryMessage, t.Image.(*image.Gray).Pix)
					case protocol.PipboyMemory:
						conn.WriteJSON(memresponse{
							Type:   "memory_update",
							Memory: t,
						})
					default:
					}
				case err := <-pc.Err():
					logger.Error(err)
					return err
				case <-localmaptimer.C:
					pc.RequestLocalMapUpdate()
					localmaptimer.Reset(75 * time.Millisecond)
				}
			}
		}()
		if err != nil {
			logger.Error(err)
			return
		}
	})

	r.Static("/assets", filepath.Join(*publicFlag, "assets"))

	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(*publicFlag, "index.html"))
	})

	go http.ListenAndServe(":8001", nil)

	r.Run(":8000")
}
