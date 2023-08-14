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
const TRIES = 6

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
	cells        [WORDS_LETTERS_COUNT]Cell
	is_filled_in bool
	saved_word   string
}

var (
	width          = 1300
	height         = 800
	word_to_guess  = ""
	cells          = [TRIES]Word{}
	attempt_index  = 0
	typed_word_str = ""
	font_face      font.Face
)
 
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	font_face, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(50),
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

		for _, cell := range cells[attempt_index].cells {
			if cell.is_empty {
				cell.character = char
				break
			}
		}
	}
	if len(typed_word_str) > 0 && inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		typed_word_str = typed_word_str[:len(typed_word_str)-1]
	}
	if len(typed_word_str) == 5 && inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		// User submits attempt
		cells[attempt_index].saved_word = typed_word_str
		cells[attempt_index].is_filled_in = true
		typed_word_str = ""
		attempt_index++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw each cell
	for row_index := 0; row_index < len(cells); row_index++ {
		var current_row = &cells[row_index]
		for i := 0; i < WORDS_LETTERS_COUNT; i++ {
			var x float32 = float32(i)*125 + 250
			var y float32 = float32(row_index) * (CELL_SIZE + 50)

			vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, color.RGBA{uint8(50 * i), 255, 255, 255}, false)
			if attempt_index == row_index {
				if i < len(typed_word_str) { // continue here
					text.Draw(screen, string(typed_word_str[i]), font_face, int(x+(CELL_SIZE/2)), int(y+(CELL_SIZE/2)), color.RGBA{255, 0, 0, 255})
				}
			} else if current_row.is_filled_in {
				text.Draw(screen, string(current_row.saved_word[i]), font_face, int(x+(CELL_SIZE/2)), int(y+(CELL_SIZE/2)), color.RGBA{255, 0, 0, 255})
			}
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
