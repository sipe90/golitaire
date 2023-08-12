package games

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sipe90/golitaire/internal/cards"
)

const (
	Suites    = 4
	Cascades  = 8
	Freecells = 4
)

type Pos struct {
	x int
	y int
}

type FreeCell struct {
	foundations []cards.Card
	freecells   []cards.Card
	cascades    [][]cards.Card
	selected    *Pos
	position    *Pos
}

func FreeCellGame() *FreeCell {
	foundations := make([]cards.Card, Suites)
	freecells := make([]cards.Card, Freecells)
	cascades := make([][]cards.Card, Cascades)

	for i := 0; i < len(cascades); i++ {
		cascades[i] = []cards.Card{}
	}

	return &FreeCell{
		foundations: foundations,
		freecells:   freecells,
		cascades:    cascades,
		selected:    nil,
		position:    &Pos{},
	}
}

func (f *FreeCell) Deal(number int) {
	deck := cards.CreateDeck()
	deck = cards.Shuffle(&deck, number)

	for i := 0; i < len(deck); i++ {
		f.cascades[i%Cascades] = append(f.cascades[i%Cascades], deck[i])
	}

	f.position.x = 0
	f.position.y = len(f.cascades[0]) - 1
}

func (f *FreeCell) Up() {
	f.position.y = int(math.Max(float64(f.position.y-1), -1))
}

func (f *FreeCell) Down() {
	f.position.y = int(math.Min(float64(f.position.y+1), float64(len(f.cascades[f.position.x])-1)))
}

func (f *FreeCell) Left() {
	f.position.x = int(math.Max(float64(f.position.x-1), 0))
	f.position.y = int(math.Min(float64(f.position.y), float64(len(f.cascades[f.position.x])-1)))
}

func (f *FreeCell) Right() {
	f.position.x = int(math.Min(float64(f.position.x+1), Cascades-1))
	f.position.y = int(math.Min(float64(f.position.y), float64(len(f.cascades[f.position.x])-1)))
}

func (f *FreeCell) Select() {}

func (f *FreeCell) View() string {
	top := lipgloss.JoinHorizontal(lipgloss.Top, f.viewFoundations(), "  ", f.viewFreeCells())
	cascades := f.viewCascades()

	return lipgloss.PlaceHorizontal(80, lipgloss.Center, top) +
		"\n" +
		lipgloss.PlaceHorizontal(80, lipgloss.Center, cascades)
}

func (f *FreeCell) viewFoundations() string {
	foundationsView := make([]string, len(f.foundations))

	for i, c := range f.foundations {
		selected := f.selected != nil && f.selected.y == -1 && f.selected.x == i
		hovered := f.position.y == -1 && f.position.x == i
		foundationsView[i] = c.View(selected, hovered, false)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, foundationsView...)
}

func (f *FreeCell) viewFreeCells() string {
	freecellsView := make([]string, len(f.freecells))

	for i, c := range f.freecells {
		selected := f.selected != nil && f.selected.y == -1 && f.selected.x == i+Suites
		hovered := f.position.y == -1 && f.position.x == i+Suites
		freecellsView[i] = c.View(selected, hovered, false)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, freecellsView...)
}

func (f *FreeCell) viewCascades() string {
	cascadesView := make([]string, len(f.cascades))

	for i, cascade := range f.cascades {
		topIdx := len(cascade) - 1

		var b strings.Builder
		b.Grow(150)

		for j, c := range cascade {
			selected := f.selected != nil && f.selected.y == j && f.selected.x == i
			hovered := f.position.y == j && f.position.x == i
			top := j == topIdx

			fmt.Fprint(&b, c.View(selected, hovered, !top))
			fmt.Fprintln(&b)
		}
		cascadesView[i] = b.String()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cascadesView...)
}

func (f *FreeCell) Debug() {
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
