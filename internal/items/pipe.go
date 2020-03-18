package items

import (
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"math/rand"
	"time"
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
	PipeImgScale float64
	GapWidth float64
	Unused bool
}

func initializePipePart(position collision2d.Vector, PipeHeadBox collision2d.Box, PipeBodyBox collision2d.Box, imgPath string) (pipePart pipePart){
	pipePart.Position = position
	pipePart.PipeHeadBox = PipeHeadBox
	pipePart.PipeBodyBox = PipeBodyBox

	var err error
	pipePart.Img,_, err = ebitenutil.NewImageFromFile(imgPath, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (pipe * Pipe) initializeUpperPipe(){
	position := collision2d.Vector{X: 500, Y: 0}
		pipe.UpperPipe = initializePipePart(
			position,
			collision2d.Box{
				Pos: collision2d.Vector{X: position.X + 0 * pipe.PipeImgScale, Y: position.Y + 775 * pipe.PipeImgScale},
				W:   (215 - 0) * pipe.PipeImgScale,
				H:   (852 - 775) * pipe.PipeImgScale,
			},
			collision2d.Box{
				Pos: collision2d.Vector{X: position.X + 31 * pipe.PipeImgScale, Y: position.Y + 0 * pipe.PipeImgScale},
				W:   (185 - 31) * pipe.PipeImgScale,
				H:   (775 - 0) * pipe.PipeImgScale,
			},
			"resources/img/upperPipe.png",
		)
	return
}

func (pipe * Pipe) initializeDownPipe() {
	position := collision2d.Vector{
		X: pipe.UpperPipe.Position.X,
		Y: pipe.GapWidth,
	}

	pipe.DownPipe = initializePipePart(
		position,
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 0 * pipe.PipeImgScale, Y: position.Y + 0 * pipe.PipeImgScale},
			W:   215 * pipe.PipeImgScale,
			H:   78 * pipe.PipeImgScale,
		},
		collision2d.Box{
			Pos: collision2d.Vector{X: position.X + 32 * pipe.PipeImgScale, Y: position.Y + 78 * pipe.PipeImgScale},
			W:   (185 - 32) * pipe.PipeImgScale,
			H:   (852 - 78) * pipe.PipeImgScale,
		},
		"resources/img/downPipe.png",
	)

	return
}

func (pipe * Pipe) Initialize(birdSize float64){
	rand.Seed(time.Now().UnixNano())

	pipe.Unused = false

	pipe.PipeImgScale = 0.5
	pipe.initializeUpperPipe()

	tolerance := 40
	min := int(pipe.UpperPipe.Position.Y +  pipe.UpperPipe.PipeBodyBox.H + pipe.UpperPipe.PipeHeadBox.H) + int(birdSize) + tolerance
	max := min + 60

	pipe.GapWidth = float64(rand.Intn(max-min+1) + min)

	pipe.initializeDownPipe()

	return
}


func (pipe * Pipe) Play(){
	pipe.UpperPipe.Position.X -= 1.8
	pipe.UpperPipe.PipeBodyBox.Pos.X -= 1.8
	pipe.UpperPipe.PipeHeadBox.Pos.X -= 1.8

	pipe.DownPipe.Position.X -= 1.8
	pipe.DownPipe.PipeBodyBox.Pos.X -= 1.8
	pipe.DownPipe.PipeHeadBox.Pos.X -= 1.8

	if pipe.UpperPipe.PipeHeadBox.Pos.X + pipe.UpperPipe.PipeHeadBox.W < 0 {
		pipe.Unused = true
	}

}

func (pipe * Pipe) Draw(screen *ebiten.Image) {
	opUpperPipe := &ebiten.DrawImageOptions{}
	opUpperPipe.GeoM.Scale(pipe.PipeImgScale, pipe.PipeImgScale)
	opUpperPipe.GeoM.Translate(pipe.UpperPipe.Position.X, pipe.UpperPipe.Position.Y)

	opDownPipe := &ebiten.DrawImageOptions{}
	opDownPipe.GeoM.Scale(pipe.PipeImgScale, pipe.PipeImgScale)
	opDownPipe.GeoM.Translate(pipe.DownPipe.Position.X, pipe.DownPipe.Position.Y)

	screen.DrawImage(pipe.UpperPipe.Img, opUpperPipe)
	screen.DrawImage(pipe.DownPipe.Img, opDownPipe)
}
