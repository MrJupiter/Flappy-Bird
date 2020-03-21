package main

import (
	"github.com/MrJupiter/Flappy-Bird/internal/ui"
	"github.com/hajimehoshi/ebiten"
	_ "image/png"
	"log"
)

func main() {
	game := new(ui.Game)
	game.Initialize()

	if err := ebiten.Run(game.Update, game.WindowDimensions.Width, game.WindowDimensions.Height, 1, "Flappy Bird"); err != nil {
		log.Fatal(err)
	}
}