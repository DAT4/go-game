package main

import (
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Players    map[byte]*Player
	Movement   map[byte]func()[]byte
	You        byte
	Conn       *websocket.Conn
	State      int
	Message    string
	MsgHistory string
}

type Sprite struct {
	left  *ebiten.Image
	right *ebiten.Image
	up    *ebiten.Image
	down  *ebiten.Image
}

type Position struct {
	x float64
	y float64
}


type Player struct {
	Sprite
	*Position
	face    *ebiten.Image
	options *ebiten.DrawImageOptions
}
