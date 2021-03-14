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



	game := setup(c)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Backend Game")

	go func() {
		for {
			_, move, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			game.Players[move[0]].options.GeoM.Apply(float64(move[1]),float64(move[2]))
		}
	}()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
