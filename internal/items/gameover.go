package items

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type GameOver struct {
	Img *ebiten.Image
}

func (gameOver *GameOver) Initialize(){
	var err error
	gameOver.Img,_, err = ebitenutil.NewImageFromFile("resources/img/gameover.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
}

func (gameOver *GameOver) GetDrawOptions(screen *ebiten.Image) *ebiten.DrawImageOptions{
	opGameOver := &ebiten.DrawImageOptions{}
	gameOverImgWidth, gameOverImgHeight := gameOver.Img.Size()
	screenWidth, screenHeight := screen.Size()
	opGameOver.GeoM.Translate(float64(screenWidth)/2 - float64(gameOverImgWidth)/2 , float64(screenHeight)/2 - float64(gameOverImgHeight)/2)

	return opGameOver
}