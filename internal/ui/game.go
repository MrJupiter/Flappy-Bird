package ui

import (
	"github.com/MrJupiter/Flappy-Bird/internal/items"
	"github.com/MrJupiter/Flappy-Bird/internal/ui/components"
	"github.com/Tarliton/collision2d"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/text"
	"gitlab.com/smartballs/driver/assets/fonts"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)



type Game struct{
	Background *components.Background
	Floor * components.Floor
	Bird *items.Bird
	Pipes[] *items.Pipe
	Score int
	WindowsDimensions dimension
	GameOver *components.GameOver
}

type dimension struct {
	Width, Height int
}

var (
	checkBirdFloorCollision, checkBirdPipeCollision bool
	fontsScore font.Face
	audioContext *audio.Context
	gameAudioPlayer  *audio.Player
	hitGameOverSound = true
)

const (
	pipesNumber  = 8
)

func (game *Game) createPipes(number int) {
	pipe := new(items.Pipe)
	birdImgHeight, _ := game.Bird.Img.Size()
	pipe.Initialize(float64(birdImgHeight)*game.Bird.ImgScale, 1024)
	game.Pipes = append(game.Pipes, pipe)

	for i:=1; i<number; i++ {
		pipeLoop := new(items.Pipe)
		pipeLoop.Initialize(float64(birdImgHeight)*game.Bird.ImgScale, game.Pipes[i-1].UpperPipe.Position.X + game.Pipes[i-1].UpperPipe.PipeHeadBox.W)
		game.Pipes = append(game.Pipes, pipeLoop)
	}
	return
}

func (game *Game) drawPipes(screen *ebiten.Image){
	for i:=0; i<len(game.Pipes) ; i++ {
		game.Pipes[i].Draw(screen)
		game.Pipes[i].Play()
	}
}


func (game *Game) getAudioPlayer(audioPath string) *audio.Player{
	f, err := os.Open(audioPath)
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.Decode(audioContext, f)
	if err != nil {
		log.Fatal(err)
	}
	audioPlayer, err := audio.NewPlayer(audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
	return audioPlayer
}


func (game *Game) initializeAudioContext(){
	if audioContext == nil{
		var err error
		audioContext, err = audio.NewContext(44100)
		if err != nil {
			log.Fatal(err)
		}
	}
	gameAudioPlayer = game.getAudioPlayer("resources/audio/gameMusic.wav")
}

func (game *Game) Initialize(){
	game.initializeAudioContext()
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	fontsScore =  truetype.NewFace(tt, &truetype.Options{
			Size:    24,
			DPI:     72,
			Hinting: font.HintingFull,
		})

	game.WindowsDimensions = dimension{Width: 1024, Height: 768}
	game.Score = 0

	game.Background = new(components.Background)
	game.Background.Initialize()

	game.Floor = new(components.Floor)
	game.Floor.Initialize()

	game.Bird = new(items.Bird)
	game.Bird.Initialize()

	game.createPipes(pipesNumber)

	game.GameOver = new(components.GameOver)
	game.GameOver.Initialize()
	return
}

func (game *Game) checkGameOverTrigger() {
	checkUpperPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].UpperPipe.PipeBodyBox.ToPolygon())
	checkUpperPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].UpperPipe.PipeHeadBox.ToPolygon())
	checkDownPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].DownPipe.PipeBodyBox.ToPolygon())
	checkDownPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].DownPipe.PipeHeadBox.ToPolygon())

	checkBirdPipeCollision = checkUpperPipeBodyBox || checkUpperPipeHeadBox || checkDownPipeBodyBox || checkDownPipeHeadBox
	checkBirdFloorCollision, _ = collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Floor.FloorBox.ToPolygon())
}

func (game *Game) incrementScore() {
	for i:=0; i<len(game.Pipes); i++{
		birdPassed, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[i].GapBox.ToPolygon())
		if birdPassed {
			game.Pipes[i].Passed = birdPassed
		}
		if game.Pipes[i].Passed && game.Bird.FlappyBox.Pos.X > game.Pipes[i].UpperPipe.PipeHeadBox.Pos.X + game.Pipes[i].UpperPipe.PipeHeadBox.W {
			game.Score++
			game.Pipes[i].Passed = false
		}
	}
}

func (game *Game) Update(screen *ebiten.Image) error {
	if !gameAudioPlayer.IsPlaying(){
		gameAudioPlayer.SetVolume(0.2)
		gameAudioPlayer.Rewind()
		gameAudioPlayer.Play()
	}

	rand.Seed(time.Now().UnixNano())

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.DrawImage(game.Background.Img, game.Background.GetDrawOptions())
	defer screen.DrawImage(game.Floor.Img, game.Floor.GetDrawOptions())
	defer text.Draw(screen, strconv.Itoa(game.Score), fontsScore, 10,30, color.Black)
	game.checkGameOverTrigger()
	if checkBirdFloorCollision || checkBirdPipeCollision {
		screen.DrawImage(game.GameOver.Img, game.GameOver.GetDrawOptions(game.WindowsDimensions.Width, game.WindowsDimensions.Height))

		if hitGameOverSound {
			hitAudioPlayer := game.getAudioPlayer("resources/audio/hit.wav")
			if !hitAudioPlayer.IsPlaying(){
				hitAudioPlayer.SetVolume(0.1)
				hitAudioPlayer.Rewind()
				hitAudioPlayer.Play()
			}
			hitGameOverSound = false
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			checkBirdFloorCollision = false
			hitGameOverSound = true
			game.Score = 0
			game.Pipes = nil
			game.Bird.Initialize()
			game.createPipes(pipesNumber)
		}
	}else{
		screen.DrawImage(game.Bird.Img, game.Bird.GetDrawOptions())

		game.drawPipes(screen)
		if game.Pipes[0].Undisplayed == true {
			game.Pipes = game.Pipes[1:]
			pipeLoop := new(items.Pipe)
			birdImgHeight, _ := game.Bird.Img.Size()
			pipeLoop.Initialize(float64(birdImgHeight)*game.Bird.ImgScale, game.Pipes[len(game.Pipes)-1].UpperPipe.Position.X + game.Pipes[len(game.Pipes)-1].UpperPipe.PipeHeadBox.W)
			game.Pipes = append(game.Pipes, pipeLoop)
		}

		game.incrementScore()

		game.Bird.Play()
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			game.Bird.Jump()
			jumpAudioPlayer := game.getAudioPlayer("resources/audio/jump.wav")
			if !jumpAudioPlayer.IsPlaying(){
				jumpAudioPlayer.SetVolume(0.1)
				jumpAudioPlayer.Rewind()
				jumpAudioPlayer.Play()
			}
		}
	}

	return nil
}