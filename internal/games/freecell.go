package games

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sipe90/golitaire/internal/cards"
)

const (
	Suites    = 4
	Cascades  = 8
	Freecells = 4

	textColor = lipgloss.Color("28")
)

var textStyle = lipgloss.NewStyle().Foreground(textColor)

// [y == -1] and [0 <= x < Suites] are foundations
// [y == -1] and [Suites <= x < Freecells - 1] are freecells
// [y >= 0] and [x >= 0] are the cascades
type Pos struct {
	x int
	y int
}

func (p *Pos) IsFoundation() bool {
	return p.x < Suites && p.y == -1
}

func (p *Pos) IsFreeCell() bool {
	return p.x >= Suites && p.y == -1
}

func (p *Pos) IsCascade() bool {
	return p.y >= 0
}

func (p *Pos) Equals(other *Pos) bool {
	return p.x == other.x && p.y == other.y
}

func (p *Pos) String() string {
	return fmt.Sprintf("[%d, %d]", p.x, p.y)
}

type FreeCell struct {
	width       int
	height      int
	number      int
	foundations []*cards.Card
	freecells   []*cards.Card
	cascades    [][]*cards.Card
	selected    *Pos
	position    *Pos
}

func FreeCellGame() *FreeCell {
	foundations := make([]*cards.Card, Suites)
	freecells := make([]*cards.Card, Freecells)
	cascades := make([][]*cards.Card, Cascades)

	for i := 0; i < len(foundations); i++ {
		foundations[i] = cards.NewCard(cards.VALUE_EMPTY, cards.Suite(i))
	}

	for i := 0; i < len(freecells); i++ {
		freecells[i] = cards.NewCard(cards.VALUE_EMPTY, cards.SUITE_EMPTY)
	}

	for i := 0; i < len(cascades); i++ {
		cascades[i] = make([]*cards.Card, 21)
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
	f.number = number
	deck := cards.CreateDeck()
	deck = cards.Shuffle(&deck, number)

	for i := 0; i < len(deck); i++ {
		f.cascades[i%Cascades][i/Cascades] = deck[i]
	}

	f.position.x = 0
	f.position.y = 6
}

func (f *FreeCell) Up() {
	f.position.y = int(math.Max(float64(f.position.y-1), -1))
}

func (f *FreeCell) Down() {
	f.position.y = int(math.Min(float64(f.position.y+1), float64(f.getTopIndex(f.position.x))))
}

func (f *FreeCell) Left() {
	f.position.x = int(math.Max(float64(f.position.x-1), 0))
	f.position.y = int(math.Min(float64(f.position.y), float64(f.getTopIndex(f.position.x))))
}

func (f *FreeCell) Right() {
	f.position.x = int(math.Min(float64(f.position.x+1), Cascades-1))
	f.position.y = int(math.Min(float64(f.position.y), float64(f.getTopIndex(f.position.x))))
}

func (f *FreeCell) Select() {
	posX := f.position.x
	posY := f.position.y

	// Select current position
	if f.selected == nil {
		if f.isCurrentPositionSelectable() {
			f.selected = &Pos{
				x: posX,
				y: posY,
			}
			log.Printf("Selected: %v", f.selected)
		}
		return
	}

	// Clear selection
	if f.selected.Equals(f.position) {
		f.selected = nil
		log.Print("Selection cleared")
		return
	}

	f.tryMoveCard(f.selected, f.position)
}

func (f *FreeCell) getTopIndex(index int) int {
	cascade := f.cascades[index]
	for i := len(cascade) - 1; i >= 0; i-- {
		if cascade[i] != nil {
			return i
		}
	}
	return 0
}

func (f *FreeCell) getCard(pos *Pos) *cards.Card {
	if pos.IsFoundation() {
		return f.foundations[pos.x]
	}
	if pos.IsFreeCell() {
		return f.freecells[pos.x-Suites]
	}
	return f.cascades[pos.x][pos.y]
}

func (f *FreeCell) tryMoveCard(from, to *Pos) {
	if f.canMoveCard(from, to) {
		f.moveCard(from, to)
	}
}

func (f *FreeCell) canMoveCard(from, to *Pos) bool {
	// Can't move onto itself
	if from.Equals(to) {
		return false
	}
	// Can't remove cards from foundations
	if from.IsFoundation() {
		return false
	}
	fromCard := f.getCard(from)
	toCard := f.getCard(to)
	// From freecell/cascade to another freecell
	if to.IsFreeCell() {
		return toCard.IsEmpty()
	}
	// From freecell/cascade to foundation
	if to.IsFoundation() {
		return cards.CanPlaceOnFoundation(fromCard, toCard)
	}
	// From freecell/cascade to cascade
	return toCard.IsEmpty() || cards.CanStack(cards.ALT_COLOR_DESC, fromCard, toCard)
}

func (f *FreeCell) moveCard(from, to *Pos) {
	if from.IsFoundation() {
		return
	}

	f.selected = nil

	card := f.getCard(from)

	log.Printf("Moving card %v from %v to %v", card, from, to)

	if to.IsFoundation() {
		f.foundations[to.x] = card
	} else if to.IsFreeCell() {
		f.freecells[to.x-Suites] = card
	} else {
		toCard := f.getCard(to)
		if to.y == 0 && toCard.IsEmpty() {
			f.cascades[to.x][to.y] = card
		} else {
			f.cascades[to.x][to.y+1] = card
			f.position.y++
		}
	}

	if from.IsFreeCell() {
		f.freecells[from.x-Suites] = cards.NewEmptyCard()
	} else {
		if from.y == 0 {
			f.cascades[from.x][from.y] = cards.NewEmptyCard()
		} else {
			f.cascades[from.x][from.y] = nil
		}
	}
}

func (f *FreeCell) View() string {
	top := lipgloss.JoinHorizontal(lipgloss.Top, f.viewFoundations(), "  ", f.viewFreeCells())
	cascades := f.viewCascades()

	gameView := textStyle.Render(fmt.Sprintf("Game #%d", f.number)) +
		"\n\n" +
		top +
		"\n" +
		cascades

	return lipgloss.PlaceHorizontal(
		f.width,
		lipgloss.Center,
		gameView,
	)
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
		var b strings.Builder
		for j, c := range cascade {
			if c == nil {
				break
			}
			selected := f.selected != nil && f.selected.y == j && f.selected.x == i
			hovered := f.position.y == j && f.position.x == i
			top := j == len(cascade)-1 || cascade[j+1] == nil

			fmt.Fprint(&b, c.View(selected, hovered, !top))
			fmt.Fprintln(&b)
		}
		cascadesView[i] = b.String()
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cascadesView...)
}

func (f *FreeCell) Resize(w int, h int) {
	f.width = w
	f.height = h
}

func (f *FreeCell) isCurrentPositionSelectable() bool {
	posX := f.position.x
	posY := f.position.y

	// Foundations are not selectable
	if f.position.IsFoundation() {
		return false
	}

	card := f.getCard(f.position)

	// Empty freecells are not selectable
	if f.position.IsFreeCell() {
		return !card.IsEmpty()
	}

	// TODO: Stack selection
	if card.IsEmpty() || posY < f.getTopIndex(posX) {
		return false
	}

	return true
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
