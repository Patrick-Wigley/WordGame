package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	// Ebiten
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	// Dataframes
	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

const WORDS_LETTERS_COUNT = 5
const CELL_SIZE = 70
const TRIES = 6

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
type KeyboardKeys struct {
	colour   color.Color
	key      string
	keyFound bool
}

var (
	width            = 700
	height           = 1000
	wordToGuess      = ""
	cells            = [TRIES]Word{}
	attempt_index    = 0
	typed_word_str   = ""
	cellsFontBIG     font.Face
	cellsFontSMALL   font.Face
	word_found       = false
	wordsDataFrame   = dataframe.DataFrame{}
	showNotvalidWord = false
	possibleChars    = []KeyboardKeys{KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "a"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "b"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "c"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "d"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "e"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "f"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "g"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "h"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "i"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "j"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "k"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "l"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "m"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "n"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "o"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "p"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "q"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "r"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "s"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "t"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "u"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "v"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "w"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "x"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "y"}, KeyboardKeys{colour: color.RGBA{255, 255, 255, 255}, key: "z"}}
	charsUsed        = []KeyboardKeys{}
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	cellsFontBIG, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(50),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	cellsFontSMALL, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(20),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	wordsFile, err := os.Open("words.csv")
	if err != nil {
		println("ERROR: ", err)
	}
	defer wordsFile.Close()
	wordsDataFrame = dataframe.ReadCSV(wordsFile)

	randomSource := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(randomSource)

	var wordsCount = wordsDataFrame.Col("WORDS").Len()
	var randomIndex = randomGenerator.Intn(wordsCount)
	wordToGuess = fmt.Sprintf("%v", wordsDataFrame.Col("WORDS").Val(randomIndex))
	println(randomIndex, wordToGuess)

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
		if showNotvalidWord {
			showNotvalidWord = false
		}
	}
	if len(typed_word_str) == 5 && inpututil.IsKeyJustReleased(ebiten.KeyEnter) {
		// User submits attempt
		if typed_word_str == wordToGuess {
			// User found word
			word_found = true
		} else {
			// check if word submitted is actually a word
			var filteredMatches = wordsDataFrame.Filter(dataframe.F{Colname: "WORDS", Comparator: series.In, Comparando: strings.ToLower(typed_word_str)})

			if filteredMatches.Col("WORDS").Len() != 0 {
				var missingCharacters = wordToGuess

				for i := 0; i < len(wordToGuess); i++ {
					var char = string(wordToGuess[i])
					var charInTypedWord = string(typed_word_str[i])

					if char == charInTypedWord {
						// Player found letter in correct place
						cells[attempt_index].cells[i].colour = color.RGBA{0, 255, 0, 255}
						missingCharacters = missingCharacters[:i] + "|" + missingCharacters[i+1:]

						var KeyboardcharIndex = charIndexInArr(char, possibleChars)
						possibleChars[KeyboardcharIndex].colour = color.RGBA{0, 255, 0, 255}
						possibleChars[KeyboardcharIndex].keyFound = true

					}
				}
				for i := 0; i < len(missingCharacters); i++ {
					var char = string(missingCharacters[i])
					var charInTypedWord = string(typed_word_str[i])

					if char != "|" {
						if strings.Contains(missingCharacters, charInTypedWord) {
							// Player found letter NOT in correct place
							cells[attempt_index].cells[i].colour = color.RGBA{255, 200, 100, 255}
							// Continue here
							var KeyboardcharIndex = charIndexInArr(charInTypedWord, possibleChars)
							if !(possibleChars[KeyboardcharIndex].keyFound) {
								possibleChars[KeyboardcharIndex].colour = color.RGBA{255, 200, 100, 255}
							}

							var charactersFirstInstance = strings.Index(missingCharacters, char)
							missingCharacters = missingCharacters[:charactersFirstInstance] + "|" + missingCharacters[charactersFirstInstance+1:]
						}
					}
				}

				cells[attempt_index].saved_word = typed_word_str
				cells[attempt_index].is_filled_in = true
				typed_word_str = ""
				attempt_index++
			} else {
				showNotvalidWord = true
				println("Not a word")
			}

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
				var x float32 = (float32(width) / float32(WORDS_LETTERS_COUNT)) + (float32(i) * (float32(CELL_SIZE) + 20.0))
				var y float32 = 100 + (float32(rowIndex) * (CELL_SIZE + 50))
				var currentCell = cells[rowIndex].cells[i]

				// Draw Cell
				if showNotvalidWord && attempt_index == rowIndex {
					vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, color.RGBA{255, 0, 0, 255}, false)
				} else {
					vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, currentCell.colour, false)
				}

				if attempt_index == rowIndex {
					if i < len(typed_word_str) {
						text.Draw(screen, string(typed_word_str[i]), cellsFontBIG, int((x + 25)), int((y + 45)), color.RGBA{0, 0, 0, 255})
					}
				} else if currentRow.is_filled_in {
					text.Draw(screen, string(currentRow.saved_word[i]), cellsFontBIG, int((x + 25)), int((y + 45)), color.RGBA{0, 0, 0, 255})
				}
			}
		}
		var keyboardLocation = [2]float32{100, 800}
		for char := 0; char < len(possibleChars); char++ {
			vector.DrawFilledRect(screen, keyboardLocation[0], keyboardLocation[1], 45, 45, possibleChars[char].colour, false)
			text.Draw(screen, possibleChars[char].key, cellsFontSMALL, int(keyboardLocation[0]+20), int(keyboardLocation[1]+25), color.RGBA{0, 0, 0, 255})
			keyboardLocation[0] += 50

			if ((char + 1) % 8) == 0 {
				keyboardLocation[1] += 50
				keyboardLocation[0] = 100
			}
		}
	} else {
		text.Draw(screen, "You Win!", cellsFontBIG, 0, (height/2)-10, color.RGBA{255, 0, 0, 255})
		text.Draw(screen, "The word was "+wordToGuess, cellsFontBIG, 0, (height/2)+40, color.RGBA{255, 0, 0, 255})

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

func charIndexInArr(value string, arr []KeyboardKeys) int {
	for i, key := range arr {
		if key.key == value {
			return i
		}
	}
	return -1
}
