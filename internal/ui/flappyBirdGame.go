package ui

import (
	"github.com/MrJupiter/Flappy-Bird/internal/items"
	"github.com/MrJupiter/Flappy-Bird/internal/ui/components"
	"github.com/MrJupiter/Flappy-Bird/resources/fonts"
	"github.com/Tarliton/collision2d"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type FlappyBirdGame struct{
	Background components.Background
	Floor components.Floor
	Bird *items.Bird
	Pipes[] *items.Pipe
	FlappyBirdScore int
	WindowDimensions dimension
	GameOver components.GameOver
}

type dimension struct {
	Width, Height int
}

var (
	fontsFlappyBirdScore font.Face
	audioContext *audio.Context
	gameAudioPlayer  *audio.Player
	hitGameOverSound = true
	startGame bool
)

const (
	pipesNumber  = 8
)

func (game *FlappyBirdGame) createPipes(number int) {
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

func (game *FlappyBirdGame) drawPipes(screen *ebiten.Image){
	for i:=0; i<len(game.Pipes) ; i++ {
		game.Pipes[i].Draw(screen)
		game.Pipes[i].Play()
	}
}


func getAudioPlayer(audioPath string) *audio.Player{
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


func initializeAudioContext(){
	if audioContext == nil{
		var err error
		audioContext, err = audio.NewContext(44100)
		if err != nil {
			log.Fatal(err)
		}
	}
	gameAudioPlayer = getAudioPlayer("resources/audio/gameMusic.wav")
}

func initializeFont(){
	tt, err := freetype.ParseFont(fonts.GetFont())
	if err != nil {
		log.Fatal(err)
	}

	fontsFlappyBirdScore =  truetype.NewFace(tt, &truetype.Options{
		Size:    34,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func (game *FlappyBirdGame) Initialize(){
	initializeAudioContext()
	initializeFont()

	game.WindowDimensions = dimension{Width: 1024, Height: 768}
	game.FlappyBirdScore = 0

	game.Background.Initialize()
	game.Floor.Initialize()
	game.GameOver.Initialize()

	game.Bird = new(items.Bird)
	game.Bird.Initialize()

	game.createPipes(pipesNumber)

	return
}

func (game *FlappyBirdGame) checkGameOverTrigger() bool {
	index := game.getIndexOfPipeThatMayCollide()
	checkUpperPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[index].UpperPipe.PipeBodyBox.ToPolygon())
	checkUpperPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[index].UpperPipe.PipeHeadBox.ToPolygon())
	checkDownPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[index].DownPipe.PipeBodyBox.ToPolygon())
	checkDownPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[index].DownPipe.PipeHeadBox.ToPolygon())

	checkBirdPipeCollision := checkUpperPipeBodyBox || checkUpperPipeHeadBox || checkDownPipeBodyBox || checkDownPipeHeadBox
	checkBirdFloorCollision, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Floor.FloorBox.ToPolygon())
	return checkBirdPipeCollision || checkBirdFloorCollision
}

func (game *FlappyBirdGame) incrementFlappyBirdScore() {
	for i:=0; i<len(game.Pipes); i++{
		birdPassed, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[i].GapBox.ToPolygon())
		if birdPassed {
			game.Pipes[i].Passed = birdPassed
		}
		if game.Pipes[i].Passed && game.Bird.FlappyBox.Pos.X > game.Pipes[i].UpperPipe.PipeHeadBox.Pos.X + game.Pipes[i].UpperPipe.PipeHeadBox.W {
			game.FlappyBirdScore++
			game.Pipes[i].Ignored = true // used to compare collision between the bird and the nearest pipe (not bypassed yet)
			game.Pipes[i].Passed = false // used to increment the score
		}
	}
}

func (game *FlappyBirdGame) getIndexOfPipeThatMayCollide() int{
	for i:=0; i<len(game.Pipes) ; i++ {
		if !game.Pipes[i].Ignored {
			return i
		}
	}
	return 0
}


func (game *FlappyBirdGame) Update(screen *ebiten.Image) error {
	rand.Seed(time.Now().UnixNano())

	if !gameAudioPlayer.IsPlaying(){
		gameAudioPlayer.SetVolume(0.2)
		gameAudioPlayer.Rewind()
		gameAudioPlayer.Play()
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	screen.DrawImage(game.Background.Img, game.Background.GetDrawOptions())
	defer screen.DrawImage(game.Floor.Img, game.Floor.GetDrawOptions())
	defer text.Draw(screen,  strconv.Itoa(game.FlappyBirdScore), fontsFlappyBirdScore, 15,40, color.White)

	if game.checkGameOverTrigger() {
		screen.DrawImage(game.GameOver.Img, game.GameOver.GetDrawOptions(game.WindowDimensions.Width, game.WindowDimensions.Height))

		if hitGameOverSound {
			hitAudioPlayer := getAudioPlayer("resources/audio/hit.wav")
			hitAudioPlayer.SetVolume(0.1)
			hitAudioPlayer.Play()
			hitGameOverSound = false
		}
		startGame = false

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			hitGameOverSound = true
			game.FlappyBirdScore = 0
			game.Pipes = nil
			game.Bird.Initialize()
			game.createPipes(pipesNumber)
		}
	}else{
		screen.DrawImage(game.Bird.Img, game.Bird.GetDrawOptions())

		if game.Pipes[0].Undisplayed == true {
			game.Pipes = game.Pipes[1:]
			pipeLoop := new(items.Pipe)
			birdImgHeight, _ := game.Bird.Img.Size()
			pipeLoop.Initialize(float64(birdImgHeight)*game.Bird.ImgScale, game.Pipes[len(game.Pipes)-1].UpperPipe.Position.X + game.Pipes[len(game.Pipes)-1].UpperPipe.PipeHeadBox.W)
			game.Pipes = append(game.Pipes, pipeLoop)
		}

		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			startGame = true
		}

		if startGame {
			game.incrementFlappyBirdScore()
			game.Bird.Play()
			game.drawPipes(screen)

			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				game.Bird.Jump()
				jumpAudioPlayer := getAudioPlayer("resources/audio/jump.wav")
				jumpAudioPlayer.SetVolume(0.1)
				jumpAudioPlayer.Play()
			}
		}
	}

	return nil
}