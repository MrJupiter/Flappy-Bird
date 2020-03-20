package components

import (
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Background struct {
	Img *ebiten.Image
	ImgScale collision2d.Vector
}

func (background *Background) Initialize(){
	background.ImgScale = collision2d.Vector{
		X: 1.7,
		Y: 1.7,
	}
	var err error
	background.Img,_, err = ebitenutil.NewImageFromFile("resources/img/background.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func (background *Background) GetDrawOptions() *ebiten.DrawImageOptions {
	opBackground := &ebiten.DrawImageOptions{}
	opBackground.GeoM.Scale(background.ImgScale.X, background.ImgScale.Y)
	return opBackground
}

