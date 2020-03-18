package items

import (
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Floor struct {
	Img *ebiten.Image
	FloorBox collision2d.Box
	ImgScale collision2d.Vector
}


func (floor *Floor) Initialize(){
	var err error
	floor.Img,_, err = ebitenutil.NewImageFromFile("resources/img/floor.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	floor.ImgScale = collision2d.Vector{
		X: 6,
		Y: 1,
	}
	_, height := floor.Img.Size()
	floor.FloorBox = collision2d.Box{
		Pos: collision2d.Vector{X: 0, Y: 768 - float64(height) + 50},
		W:   336 * floor.ImgScale.X,
		H:   112 * floor.ImgScale.Y,
	}
}

func (floor *Floor) GetDrawOptions() *ebiten.DrawImageOptions{
	opFloor := &ebiten.DrawImageOptions{}
	opFloor.GeoM.Scale(floor.ImgScale.X, floor.ImgScale.Y)
	opFloor.GeoM.Translate(floor.FloorBox.Pos.X, floor.FloorBox.Pos.Y)
	return opFloor
}
