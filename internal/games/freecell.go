package games

import (
	"fmt"
	"math"

	"github.com/sipe90/golitaire/internal/cards"
)

const (
	Suites    = 4
	Cascades  = 8
	Freecells = 4
)

type Loc struct {
	x int
	y int
}

type FreeCell struct {
	foundations []cards.Card
	freecells   []cards.Card
	cascades    [][]cards.Card
	selected    Loc
	position    Loc
}

func FreeCellGame() FreeCell {
	foundations := make([]cards.Card, Suites)
	freecells := make([]cards.Card, Freecells)
	cascades := make([][]cards.Card, Cascades)

	for i := 0; i < len(cascades); i++ {
		cascades[i] = []cards.Card{}
	}

	return FreeCell{
		foundations: foundations,
		freecells:   freecells,
		cascades:    cascades,
	}
}

func (f FreeCell) Deal(number int) {
	deck := cards.CreateDeck()
	deck = cards.Shuffle(&deck, number)

	for i := 0; i < len(deck); i++ {
		f.cascades[i%Cascades] = append(f.cascades[i%Cascades], deck[i])
	}
}

func (f FreeCell) Up() {}

func (f FreeCell) Down() {}

func (f FreeCell) Left() {}

func (f FreeCell) Right() {}

func (f FreeCell) Select() {}

func (f FreeCell) Debug() {
	for _, c := range f.foundations {
		if c.IsEmpty() {
			fmt.Printf("%4v", "X")
		} else {
			fmt.Printf("%4v", c)
		}
	}
	for _, c := range f.freecells {
		if c.IsEmpty() {
			fmt.Printf("%4v", "X")
		} else {
			fmt.Printf("%4v", c)
		}
	}
	fmt.Printf("\n\n")
	var maxLen = 0
	for _, c := range f.cascades {
		maxLen = int(math.Max(float64(maxLen), float64(len(c))))
	}
	for y := 0; y < maxLen; y++ {
		for x := 0; x < len(f.cascades); x++ {
			if len(f.cascades[x]) < y+1 {
				fmt.Print("    ")
			} else {
				c := f.cascades[x][y]
				fmt.Printf("%4v", c)
			}
		}
		fmt.Println()
	}
}
