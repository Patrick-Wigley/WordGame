package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	font_face      font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	font_face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(20),
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendJustReleasedKeys(g.keys[:0])

	for _, val := range g.keys {
		char := ebiten.KeyName(val)

		if len(typed_word_str) < 5 {
			typed_word_str += char
		}

		for _, cell := range word.cells {
			if cell.is_empty {
				cell.character = char
				break
			}
		}
	}
	if len(typed_word_str) > 0 && inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		typed_word_str = typed_word_str[:len(typed_word_str)-1]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each cell
	for i := 0; i < WORDS_LETTERS_COUNT; i++ {
		var x float32 = float32(i)*125 + 250
		var y float32 = 50.0

		vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, color.RGBA{uint8(50 * i), 255, 255, 255}, false)
		if i < len(typed_word_str) { // continue here
			text.Draw(screen, string(typed_word_str[i]), font_face, int(x), int(y), color.RGBA{255, 0, 0, 0})
		}
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
