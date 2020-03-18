package ui

import (
	"github.com/MrJupiter/flappy_bird/internal/items"
	"github.com/Tarliton/collision2d"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"os"

	"log"
)



type Game struct{
	Bird *items.Bird
	Pipes[] *items.Pipe
	Floor *items.Floor
	Background *items.Background
	Score int
}

var (
	gameOver *items.GameOver
	audioContext *audio.Context
	gameAudioPlayer  *audio.Player
	jumpAudioPlayer  *audio.Player
	hitAudioPlayer  *audio.Player
	counter int
)

func (game *Game) Initialize(){
	game.Background = new(items.Background)
	game.Background.Initialize()

	game.Floor = new(items.Floor)
	game.Floor.Initialize()

	game.Bird = new(items.Bird)
	game.Bird.Initialize()

	pipe := new(items.Pipe)
	pipe.Initialize(game.Bird.FlappyBox.H)
	game.Pipes = append(game.Pipes, pipe)

	if audioContext == nil{
		var err error
		audioContext, err = audio.NewContext(44100)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.Open("resources/audio/gameMusic.wav")
		if err != nil {
			log.Fatal(err)
		}
		d, err := wav.Decode(audioContext, f)
		if err != nil {
			log.Fatal(err)
		}
		gameAudioPlayer, err = audio.NewPlayer(audioContext, d)
		if err != nil {
			log.Fatal(err)
		}
	}

	game.Score = 0

	gameOver = new(items.GameOver)
	gameOver.Initialize()
	return
}

func (game *Game) checkPipesCollision() bool{
	checkUpperPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].UpperPipe.PipeBodyBox.ToPolygon())
	checkUpperPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].UpperPipe.PipeHeadBox.ToPolygon())

	checkDownPipeBodyBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].DownPipe.PipeBodyBox.ToPolygon())
	checkDownPipeHeadBox, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Pipes[0].DownPipe.PipeHeadBox.ToPolygon())

	checkFloorCollision, _ := collision2d.TestPolygonPolygon(game.Bird.FlappyBox.ToPolygon(), game.Floor.FloorBox.ToPolygon())

	return checkDownPipeBodyBox || checkDownPipeHeadBox || checkUpperPipeHeadBox || checkUpperPipeBodyBox  || checkFloorCollision
}

func (game *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	if !gameAudioPlayer.IsPlaying(){
		gameAudioPlayer.SetVolume(0.2)
		gameAudioPlayer.Rewind()
		gameAudioPlayer.Play()
	}

	screen.DrawImage(game.Background.Img, game.Background.GetDrawOptions())

	if game.checkPipesCollision() {
		screen.DrawImage(gameOver.Img, gameOver.GetDrawOptions(screen))
		counter++
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			game.Pipes = game.Pipes[1:]
			pipe := new(items.Pipe)
			pipe.Initialize(game.Bird.FlappyBox.H)
			game.Pipes = append(game.Pipes, pipe)
			game.Initialize()
			counter = 0
		}
		if counter == 1 {
			f, err := os.Open("resources/audio/hit.wav")
			if err != nil {
				log.Fatal(err)
			}
			d, err := wav.Decode(audioContext, f)
			if err != nil {
				log.Fatal(err)
			}
			jumpAudioPlayer, err = audio.NewPlayer(audioContext, d)
			if err != nil {
				log.Fatal(err)
			}
			if !jumpAudioPlayer.IsPlaying() {
				jumpAudioPlayer.SetVolume(0.1)
				jumpAudioPlayer.Rewind()
				jumpAudioPlayer.Play()
			}
		}
	}else{
		if len(game.Pipes) > 0 {
			if game.Pipes[0].Unused == false {
				game.Pipes[0].Draw(screen)
				game.Pipes[0].Play()
			} else {
				game.Pipes = game.Pipes[1:]
				pipe := new(items.Pipe)
				pipe.Initialize(game.Bird.FlappyBox.H)
				game.Pipes = append(game.Pipes, pipe)
			}
		}

		opBird := &ebiten.DrawImageOptions{}
		opBird.GeoM.Scale(game.Bird.ImgScale, game.Bird.ImgScale)
		opBird.GeoM.Translate(game.Bird.Position.X, game.Bird.Position.Y)
		screen.DrawImage(game.Bird.Img, opBird)

		game.Bird.Play()

		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			game.Bird.Jump()
			f, err := os.Open("resources/audio/jump.wav")
			if err != nil {
				log.Fatal(err)
			}
			d, err := wav.Decode(audioContext, f)
			if err != nil {
				log.Fatal(err)
			}
			jumpAudioPlayer, err = audio.NewPlayer(audioContext, d)
			if err != nil {
				log.Fatal(err)
			}
			if !jumpAudioPlayer.IsPlaying(){
				jumpAudioPlayer.SetVolume(0.1)
				jumpAudioPlayer.Rewind()
				jumpAudioPlayer.Play()
			}
		}
	}

	screen.DrawImage(game.Floor.Img, game.Floor.GetDrawOptions())
	return nil
}