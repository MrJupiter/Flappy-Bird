package items

import (
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type Bird struct{
	Position  collision2d.Vector
	ImgScale  float64
	Img *ebiten.Image
	FlappyBox collision2d.Box
}

func (bird * Bird) Initialize(){
	bird.Position.X, bird.Position.Y, bird.ImgScale = 25, 768/2, 0.1
	bird.FlappyBox = collision2d.Box{
		Pos: collision2d.Vector{X:bird.Position.X + 191 * bird.ImgScale, Y:bird.Position.Y + 375 * bird.ImgScale},
		W:   670 * bird.ImgScale,
		H:   490 * bird.ImgScale,
	}

	var err error
	bird.Img,_, err = ebitenutil.NewImageFromFile("resources/img/bird.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (bird *Bird) Jump(){
	if bird.Position.Y + bird.FlappyBox.Pos.Y> 0 {
		bird.Position.Y -= 6
		bird.FlappyBox.Pos.Y -= 6
	}
	return
}

func (bird *Bird) Play(){
	bird.Position.Y++
	bird.FlappyBox.Pos.Y++
	return
}