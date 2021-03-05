package main

import (
	"fmt"
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

	game := setup(c)

	go func() {
		for {
			_, move, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println("RECV:", string(move))
		}
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
