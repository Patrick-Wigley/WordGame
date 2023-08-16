package main

import (
	"image/color"
	"log"
	"strings"

	//"unicode"

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
const WORD_TO_GUESS = "hello"

type Game struct {
	keys []ebiten.Key
}

type Cell struct {
	character    string
	position     int32
	string_typed string
	is_empty     bool
	colour       color.Color
	isFound      bool
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
	word_found     = false
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

	// Initialise cells
	for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
		for cellIndex := 0; cellIndex < WORDS_LETTERS_COUNT; cellIndex++ {
			// Cell Colour
			cells[rowIndex].cells[cellIndex].colour = color.RGBA{255, 255, 255, 255}

		}
	}
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
		if typed_word_str == WORD_TO_GUESS {
			// User found word
			word_found = true
		} else {

			//var characters_found_maxed = ""
			var missing_characters = WORD_TO_GUESS

			println(attempt_index)
			for i := 0; i < len(WORD_TO_GUESS); i++ {
				var char = string(WORD_TO_GUESS[i])
				var char_in_typed_word = string(typed_word_str[i])

				if char == char_in_typed_word {
					// Player found letter in correct place
					//characters_to_processed += strings.ToUpper(char)
					cells[attempt_index].cells[i].colour = color.RGBA{0, 255, 0, 255}

					//var characters_first_instance = strings.Index(missing_characters, char)
					missing_characters = missing_characters[:i] + "|" + missing_characters[i+1:]
					println(missing_characters)
				}

			}
			for i := 0; i < len(missing_characters); i++ {
				var char = string(missing_characters[i])
				var char_in_typed_word = string(typed_word_str[i])

				if char != "|" {
					if strings.Contains(missing_characters, char_in_typed_word) {
						cells[attempt_index].cells[i].colour = color.RGBA{255, 200, 100, 255}
						var charactersFirstInstance = strings.Index(missing_characters, char)
						missing_characters = missing_characters[:charactersFirstInstance] + "|" + missing_characters[charactersFirstInstance+1:]
					}
				}
			}

			// for i := 0; i < len(typed_word_str); i++ {
			// 	var char = string(typed_word_str[i])
			// 	if strings.Contains(missing_characters, char) {
			// 		// If char in missing characters
			// 		if !(cells[attempt_index].cells[i].isFound) {
			// 			cells[attempt_index].cells[i].colour = color.RGBA{255, 200, 100, 255}
			// 		}

			// 		var characters_first_instance = strings.Index(missing_characters, char)
			// 		missing_characters = missing_characters[:characters_first_instance] + "|" + missing_characters[characters_first_instance+1:]
			// 		println(missing_characters)
			// 	}

			// }

			// else if strings.Contains(WORD_TO_GUESS, char_in_typed_word) {
			// 	// This char is in the word, not in this current location
			// 	if !(strings.Contains(characters_to_processed, char_in_typed_word)) {
			// 		// Add character to be processed
			// 		characters_to_processed += char_in_typed_word
			// 	}
			// }

			//

			// for i := 0; i < len(typed_word_str); i++ {
			// 	var char = string(typed_word_str[i])
			// 	if strings.Contains(WORD_TO_GUESS, char) {
			// 		var occurances_in_typed_word = strings.Count(typed_word_str, char)
			// 		var occurances_count = strings.Count(WORD_TO_GUESS, char)

			// 		if occurances_in_typed_word > occurances_count {
			// 			// used char more times than needed

			// 			for char_occurance := 0; char_occurance < occurances_count; char_occurance++ {

			// 			}

			// 		}

			// 		// Character found but NOT in correct place
			// 		if strings.Count(characters_processed, char) >= strings.Count(WORD_TO_GUESS, char) {
			// 			//ignore
			// 			println("Used letter more times than is in string")
			// 		} else {
			// 			cells[attempt_index].cells[i].colour = color.RGBA{255, 200, 100, 255}
			// 			println("Found letter but not in correct place")
			// 		}
			// 		if string(WORD_TO_GUESS[i]) == char {
			// 			// Character found in correct
			// 			cells[attempt_index].cells[i].colour = color.RGBA{0, 255, 0, 255}
			// 			println("Found letter in exact place")
			// 		}

			// 		// if strings.Count(typed_word_str, char) > strings.Count(WORD_TO_GUESS, char) {
			// 		// 	// Ignore
			// 	}
			// 	characters_processed += char
			// }

			cells[attempt_index].saved_word = typed_word_str
			cells[attempt_index].is_filled_in = true
			typed_word_str = ""
			attempt_index++
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if !word_found {
		// Draw each cell
		for rowIndex := 0; rowIndex < len(cells); rowIndex++ {
			var currentRow = &cells[rowIndex]
			for i := 0; i < WORDS_LETTERS_COUNT; i++ {
				var x float32 = float32(i)*125 + 250
				var y float32 = float32(rowIndex) * (CELL_SIZE + 50)
				var currentCell = cells[rowIndex].cells[i]

				vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, currentCell.colour, false)
				if attempt_index == rowIndex {
					if i < len(typed_word_str) { // continue here
						text.Draw(screen, string(typed_word_str[i]), font_face, int(x+(CELL_SIZE/2)), int(y+(CELL_SIZE/2)), color.RGBA{255, 0, 0, 255})
					}
				} else if currentRow.is_filled_in {
					text.Draw(screen, string(currentRow.saved_word[i]), font_face, int(x+(CELL_SIZE/2)), int(y+(CELL_SIZE/2)), color.RGBA{255, 0, 0, 255})
				}
			}
		}
	} else {
		text.Draw(screen, "You Win, The word was "+WORD_TO_GUESS, font_face, (width/2)-140, (height/2)-10, color.RGBA{255, 0, 0, 255})
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
