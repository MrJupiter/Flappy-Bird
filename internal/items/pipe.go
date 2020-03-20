package items

import (
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math/rand"
)

type pipePart struct{
	Position collision2d.Vector
	Img *ebiten.Image
	PipeHeadBox collision2d.Box
	PipeBodyBox collision2d.Box
}

type Pipe struct {
	DownPipe pipePart
	UpperPipe pipePart
	GapBox collision2d.Box
	Undisplayed bool
	PipeImgScale collision2d.Vector
	Passed bool
}

func randFloat(min, max float64) float64{
	return min + rand.Float64() * (max - min)
}

func initializePipePart(position collision2d.Vector, PipeBodyBox collision2d.Box, PipeHeadBox collision2d.Box, imgPath string) (pipePart pipePart){
	pipePart.Position = position
	pipePart.PipeBodyBox = PipeBodyBox
	pipePart.PipeHeadBox = PipeHeadBox

	var err error
	pipePart.Img,_, err = ebitenutil.NewImageFromFile(imgPath, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pipe * Pipe) initializeUpperPipe(pipeStartMinPositionX float64){
	pipeStartPositionX := randFloat(pipeStartMinPositionX+60, pipeStartMinPositionX + 320)

	randomY := randFloat(-150,5)
	position := collision2d.Vector{X: pipeStartPositionX, Y: randomY}
	pipe.UpperPipe = initializePipePart(
		position,
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 31 * pipe.PipeImgScale.X, Y: position.Y + 0 * pipe.PipeImgScale.Y},
			W:   (185 - 31) * pipe.PipeImgScale.X,
			H:   (775 - 0) * pipe.PipeImgScale.Y,
		},
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 0 * pipe.PipeImgScale.X, Y: position.Y + 775 * pipe.PipeImgScale.Y},
			W:   (215 - 0) * pipe.PipeImgScale.X,
			H:   (852 - 775) * pipe.PipeImgScale.Y,
		},
		"resources/img/upperPipe.png",
	)
	return
}

func (pipe * Pipe) initializeDownPipe() {
	position := collision2d.Vector{
		X: pipe.UpperPipe.Position.X,
		Y: pipe.GapBox.H + pipe.GapBox.Pos.Y,
	}

	pipe.DownPipe = initializePipePart(
		position,
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 32 * pipe.PipeImgScale.X, Y: position.Y + 78 * pipe.PipeImgScale.Y},
			W:   (185 - 32) * pipe.PipeImgScale.X,
			H:   (852 - 78) * pipe.PipeImgScale.Y,
		},
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 0 * pipe.PipeImgScale.X, Y: position.Y + 0 * pipe.PipeImgScale.Y},
			W:   215 * pipe.PipeImgScale.X,
			H:   78 * pipe.PipeImgScale.Y,
		},
		"resources/img/downPipe.png",
	)

	return
}

func (pipe * Pipe) Initialize(birdImgSize float64, pipeStartPositionX float64){
	pipe.Undisplayed = false
	pipe.Passed = false
	pipe.PipeImgScale = collision2d.Vector{
		X: 0.5,
		Y: 0.5,
	}
	pipe.initializeUpperPipe(pipeStartPositionX)

	min := birdImgSize
	gapHeight := randFloat(min, min + 50)

	pipe.GapBox = collision2d.Box{
		Pos: collision2d.Vector{X: pipe.UpperPipe.PipeHeadBox.Pos.X, Y: pipe.UpperPipe.Position.Y +  pipe.UpperPipe.PipeBodyBox.H + pipe.UpperPipe.PipeHeadBox.H },
		W:   pipe.UpperPipe.PipeHeadBox.W,
		H:   gapHeight,
	}

	pipe.initializeDownPipe()

	return
}

func (pipe * Pipe) Play(){
	pipe.UpperPipe.Position.X -= 2
	pipe.UpperPipe.PipeBodyBox.Pos.X -= 2
	pipe.UpperPipe.PipeHeadBox.Pos.X -= 2

	pipe.GapBox.Pos.X -= 2

	pipe.DownPipe.Position.X -= 2
	pipe.DownPipe.PipeBodyBox.Pos.X -= 2
	pipe.DownPipe.PipeHeadBox.Pos.X -= 2

	if pipe.UpperPipe.PipeHeadBox.Pos.X + pipe.UpperPipe.PipeHeadBox.W < 0 {
		pipe.Undisplayed = true
	}
}

func (pipe * Pipe) Draw(screen *ebiten.Image) {
	opUpperPipe := &ebiten.DrawImageOptions{}
	opUpperPipe.GeoM.Scale(pipe.PipeImgScale.X, pipe.PipeImgScale.Y)
	opUpperPipe.GeoM.Translate(pipe.UpperPipe.Position.X, pipe.UpperPipe.Position.Y)

	opDownPipe := &ebiten.DrawImageOptions{}
	opDownPipe.GeoM.Scale(pipe.PipeImgScale.X, pipe.PipeImgScale.Y)
	opDownPipe.GeoM.Translate(pipe.DownPipe.Position.X, pipe.DownPipe.Position.Y)

	screen.DrawImage(pipe.UpperPipe.Img, opUpperPipe)
	screen.DrawImage(pipe.DownPipe.Img, opDownPipe)
}