package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const WORDS_LETTERS_COUNT = 5
const CELL_SIZE = 100

type Game struct {
	keys []ebiten.Key
}

type Cell struct {
	character    string
	position     int32
	string_typed string
	is_empty     bool
}
type Word struct {
	cells [WORDS_LETTERS_COUNT]Cell
}

var (
	width          = 1300
	height         = 800
	word_to_guess  = ""
	word           = Word{}
	typed_word_str = ""
)

func (g *Game) Update() error {
	g.keys = inpututil.AppendJustReleasedKeys(g.keys[:0])

	for _, val := range g.keys {
		char := ebiten.KeyName(val)

		typed_word_str += char
		for _, cell := range word.cells {
			if cell.is_empty {
				cell.character = char
				break
			}
		}
	}

	println(typed_word_str)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each cell
	for i := 0; i < WORDS_LETTERS_COUNT; i++ {
		vector.DrawFilledRect(screen, float32(i)*125+250, 50, CELL_SIZE, CELL_SIZE, color.RGBA{uint8(50 * i), 255, 255, 255}, false)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return width, height
}

func main() {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Word Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
