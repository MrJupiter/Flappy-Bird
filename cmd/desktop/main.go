package main

import (
	"github.com/MrJupiter/flappy_bird/internal/ui"
	"github.com/hajimehoshi/ebiten"
	_ "image/png"
	"log"
)

func main() {
	game := new(ui.Game)
	game.Initialize()

	if err := ebiten.Run(game.Update, 1024, 768, 1, "Flappy Bird"); err != nil {
		log.Fatal(err)
	}
}