package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)



func main() {
	c, err := setupConnection()
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Backend Game")

	channel := make(chan []byte)

	game := setup(c,channel)

	go func() {
		for {
			_, move, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			channel <- move
		}
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
