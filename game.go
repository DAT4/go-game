package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"golang.org/x/image/colornames"
)

func (g *Game) MoveActualPlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		_ = g.Conn.WriteMessage(websocket.TextMessage, g.Movement[LEFT]())
	} else if ebiten.IsKeyPressed(ebiten.KeyJ) {
		_ = g.Conn.WriteMessage(websocket.TextMessage, g.Movement[DOWN]())
	} else if ebiten.IsKeyPressed(ebiten.KeyK) {
		_ = g.Conn.WriteMessage(websocket.TextMessage, g.Movement[UP]())
	} else if ebiten.IsKeyPressed(ebiten.KeyL) {
		_ = g.Conn.WriteMessage(websocket.TextMessage, g.Movement[RIGHT]())
	} else if repeatingKeyPressed(ebiten.KeySemicolon) {
		g.State = TYPING
	}
}

func (g *Game) writeChat() {
	g.Message += string(ebiten.InputChars())
	if repeatingKeyPressed(ebiten.KeyEnter) {
		g.State = PLAYING
		err := g.Conn.WriteMessage(websocket.TextMessage, []byte("msg: "+g.Message))
		g.Message = ""
		if err != nil {
			fmt.Println(err)
		}
	}
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			if len(g.Message) >= 1 {
				g.Message = g.Message[:len(g.Message)-1]
			}
		}
	}
}

func (g *Game) Update() error {
	if g.State == 0 {
		g.MoveActualPlayer()
	} else {
		g.writeChat()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Green)
	for _, player := range g.Players {
		screen.DrawImage(player.face,player.options)
	}
	ebitenutil.DebugPrint(screen, g.MsgHistory)
	if g.State == TYPING {
		ebitenutil.DebugPrintAt(screen, g.Message, 20, screenHeight-40)
	}
}
