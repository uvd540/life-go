package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	gameSize    int = 600
	displaySize int = 800
)

var (
	aliveColor               = color.White
	deadColor                = color.Black
	oneLivingNeighborColor   = color.RGBA{255, 0, 0, 255}
	twoLivingNeighborColor   = color.RGBA{0, 255, 0, 255}
	threeLivingNeighborColor = color.RGBA{0, 0, 255, 255}
)

type Game struct {
	buffer0        [gameSize * gameSize]bool
	buffer1        [gameSize * gameSize]bool
	currentBuffer  *[gameSize * gameSize]bool
	previousBuffer *[gameSize * gameSize]bool
	image          *ebiten.Image
	isBuffer0      bool
}

func GetAtIndex(buffer *[gameSize * gameSize]bool, x, y int) bool {
	return buffer[x*gameSize+y]
}

func SetAtIndex(buffer *[gameSize * gameSize]bool, x, y int, value bool) {
	buffer[x*gameSize+y] = value
}

func GetNumLivingNeighbors(buffer *[gameSize * gameSize]bool, x, y int) int {
	var numLivingNeighbors int = 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if !(i == 0 && j == 0) && GetAtIndex(buffer, x+i, y+j) {
				numLivingNeighbors++
			}
		}
	}
	return numLivingNeighbors
}

func (game *Game) Init() {
	game.image.Fill(deadColor)
	for x := 1; x < gameSize-1; x++ {
		for y := 1; y < gameSize-1; y++ {
			if rand.Intn(2) == 1 {
				SetAtIndex(&game.buffer1, x, y, true)
				game.image.Set(x, y, aliveColor)
			}
		}
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return gameSize, gameSize
}

func (game *Game) Update() error {
	if game.isBuffer0 {
		game.currentBuffer = &game.buffer1
		game.previousBuffer = &game.buffer0
	} else {
		game.currentBuffer = &game.buffer0
		game.previousBuffer = &game.buffer1
	}
	for x := 1; x < gameSize-1; x++ {
		for y := 1; y < gameSize-1; y++ {
			numNeighbors := GetNumLivingNeighbors(game.previousBuffer, x, y)
			if GetAtIndex(game.previousBuffer, x, y) { // live cells
				if numNeighbors < 2 || numNeighbors > 3 {
					SetAtIndex(game.currentBuffer, x, y, false)
					game.image.Set(x, y, deadColor)
				} else {
					SetAtIndex(game.currentBuffer, x, y, true)
				}
			} else { // dead cells
				if numNeighbors == 3 {
					SetAtIndex(game.currentBuffer, x, y, true)
					game.image.Set(x, y, aliveColor)
				} else {
					SetAtIndex(game.currentBuffer, x, y, false)
				}
			}
		}
	}
	game.isBuffer0 = !game.isBuffer0
	// }
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(game.image, nil)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%v", ebiten.ActualTPS()))
}

func main() {
	ebiten.SetWindowSize(displaySize, displaySize)
	ebiten.SetWindowTitle("life")
	game := &Game{}
	game.image = ebiten.NewImage(gameSize, gameSize)
	game.Init()
	if err := ebiten.RunGame(game); err != nil {
		log.Println(err)
	}
}
