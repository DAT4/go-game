package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/png"
	"log"
)

func (g *Game) setupMainPlayerMovement() map[byte]func()[]byte {
	return map[byte]func()[]byte{
		LEFT: func() []byte {
			x := g.Players[g.You].x - 1
			y := g.Players[g.You].y
			return []byte{g.You,byte(x),byte(y)}
		},
		RIGHT: func()[]byte {
			x := g.Players[g.You].x + 1
			y := g.Players[g.You].y
			return []byte{g.You,byte(x),byte(y)}
		},
		UP: func()[]byte {
			x := g.Players[g.You].x
			y := g.Players[g.You].y - 1
			return []byte{g.You,byte(x),byte(y)}
		},
		DOWN: func()[]byte {
			x := g.Players[g.You].x
			y := g.Players[g.You].y + 1
			return []byte{g.You,byte(x),byte(y)}
		},
	}
}

func (g *Game) setupPlayerSprite(playerId byte) Sprite {
	return Sprite{
		left:  getImg("images/player" + string(playerId) + "_l.png"),
		right: getImg("images/player" + string(playerId) + "_r.png"),
		up:    getImg("images/player" + string(playerId) + "_u.png"),
		down:  getImg("images/player" + string(playerId) + "_d.png"),
	}
}

func setup(c *websocket.Conn, channel <- chan []byte) (g *Game) {
	g = &Game{
		Conn: c,
		Channel: channel,
	}

	type message struct {
		command  byte
		playerId byte
		startPos Position
	}

	const (
		READY = iota
		CREATE
		ASSIGN
	)

	ok := make(chan bool)

	//Hashmap of functions which are gonna be called from the server
	commands := map[byte]func(message) error{
		READY: func(msg message) (err error) {
			ok <- true
			return
		},
		ASSIGN: func(msg message) (err error) {
			g.You = msg.playerId
			g.Movement = g.setupMainPlayerMovement()
			return
		},
		CREATE: func(msg message) (err error) {
			g.Players[msg.playerId] = &Player{
				options: &ebiten.DrawImageOptions{
					GeoM: ebiten.GeoM{},
				},
				Sprite:   g.setupPlayerSprite(msg.playerId),
			}
			g.Players[msg.playerId].face = g.Players[msg.playerId].up
			return
		},
	}

	//Background loop waiting for the server to be ready
	go func() {
		var msg message
		for {
			err := g.Conn.ReadJSON(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			f, ok := commands[msg.command]
			if ok {
				err = f(msg)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	fmt.Println("Waiting for other players...")
	<-ok
	return g
}

func getImg(path string) *ebiten.Image {
	file, err := ebitenutil.OpenFile(path)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	out := ebiten.NewImageFromImage(img)
	return out
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
