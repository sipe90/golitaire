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
	Suites           = 4
	Cascades         = 8
	Freecells        = 4
	MaxCascadeLength = 21

	textColor = lipgloss.Color("28")
)

var textStyle = lipgloss.NewStyle().Foreground(textColor)

// [y == -1] and [0 <= x < Suites] are foundations
func IsFoundation(p *cards.Pos) bool {
	return p.X < Suites && p.Y == -1
}

// [y == -1] and [Suites <= x < Freecells - 1] are freecells
func IsFreeCell(p *cards.Pos) bool {
	return p.X >= Suites && p.Y == -1
}

// [y >= 0] and [x >= 0] are the cascades
func IsCascade(p *cards.Pos) bool {
	return p.Y >= 0
}

type FreeCell struct {
	width       int
	height      int
	number      int
	foundations []*cards.Card
	freecells   []*cards.Card
	cascades    [][]*cards.Card
	selected    *cards.Pos
	position    *cards.Pos
	moves       *cards.MoveHistory
}

func FreeCellGame() *FreeCell {
	foundations := make([]*cards.Card, Suites)
	freecells := make([]*cards.Card, Freecells)
	cascades := make([][]*cards.Card, Cascades)

	return &FreeCell{
		foundations: foundations,
		freecells:   freecells,
		cascades:    cascades,
		selected:    nil,
		position:    &cards.Pos{},
		moves:       &cards.MoveHistory{},
	}
}

func (f *FreeCell) Deal(number int) {
	f.number = number

	for i := 0; i < len(f.foundations); i++ {
		f.foundations[i] = cards.NewCard(cards.VALUE_EMPTY, cards.Suite(i))
	}

	for i := 0; i < len(f.freecells); i++ {
		f.freecells[i] = cards.NewCard(cards.VALUE_EMPTY, cards.SUITE_EMPTY)
	}

	for i := 0; i < len(f.cascades); i++ {
		f.cascades[i] = make([]*cards.Card, MaxCascadeLength)
	}

	deck := cards.CreateDeck()
	deck = cards.Shuffle(&deck, number)

	for i := 0; i < len(deck); i++ {
		f.cascades[i%Cascades][i/Cascades] = deck[i]
	}

	f.position.X = 0
	f.position.Y = 6

	log.Printf("Dealt game #%v", f.number)
}

func (f *FreeCell) Redeal() {
	f.Deal(f.number)
}

func (f *FreeCell) Up() {
	f.position.Y = int(math.Max(float64(f.position.Y-1), -1))
}

func (f *FreeCell) Down() {
	f.position.Y = int(math.Min(float64(f.position.Y+1), float64(f.getTopIndex(f.position.X))))
}

func (f *FreeCell) Left() {
	f.position.X = int(math.Max(float64(f.position.X-1), 0))
	f.position.Y = int(math.Min(float64(f.position.Y), float64(f.getTopIndex(f.position.X))))
}

func (f *FreeCell) Right() {
	f.position.X = int(math.Min(float64(f.position.X+1), Cascades-1))
	f.position.Y = int(math.Min(float64(f.position.Y), float64(f.getTopIndex(f.position.X))))
}

func (f *FreeCell) Select() {
	posX := f.position.X
	posY := f.position.Y

	// Select current position
	if f.selected == nil {
		if f.isCurrentPositionSelectable() {
			f.selected = &cards.Pos{
				X: posX,
				Y: posY,
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

func (f *FreeCell) Undo() {
	if f.moves.Size() > 0 {
		move := f.moves.Pop()
		f.moveCard(&move.To, &move.From)
	}
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

func (f *FreeCell) getCard(pos *cards.Pos) *cards.Card {
	if IsFoundation(pos) {
		return f.foundations[pos.X]
	}
	if IsFreeCell(pos) {
		return f.freecells[pos.X-Suites]
	}
	return f.cascades[pos.X][pos.Y]
}

func (f *FreeCell) tryMoveCard(from, onto *cards.Pos) {
	if f.canMoveCard(from, onto) {
		var to *cards.Pos
		if IsCascade(onto) {
			to = onto.Add(0, 1)
		} else {
			to = onto
		}

		f.moveCard(from, to)
		f.moves.Push(&cards.Move{From: *from, To: *to})
		log.Printf("Moves: %v", f.moves)
	}
}

func (f *FreeCell) canMoveCard(from, onto *cards.Pos) bool {
	// Can't move onto itself
	if from.Equals(onto) {
		return false
	}
	// Can't remove cards from foundations
	if IsFoundation(from) {
		return false
	}
	fromCard := f.getCard(from)
	toCard := f.getCard(onto)
	// From freecell/cascade to another freecell
	if IsFreeCell(onto) {
		return toCard.IsEmpty()
	}
	// From freecell/cascade to foundation
	if IsFoundation(onto) {
		return cards.CanPlaceOnFoundation(fromCard, toCard)
	}
	// From freecell/cascade to cascade
	return toCard.IsEmpty() || cards.CanStack(cards.ALT_COLOR_DESC, fromCard, toCard)
}

func (f *FreeCell) moveCard(from, to *cards.Pos) {
	f.selected = nil

	card := f.getCard(from)

	log.Printf("Moving card %v from %v to %v", card, from, to)

	if IsFoundation(to) {
		f.foundations[to.X] = card
	} else if IsFreeCell(to) {
		f.freecells[to.X-Suites] = card
	} else {
		f.cascades[to.X][to.Y] = card
	}

	if IsFreeCell(from) {
		f.freecells[from.X-Suites] = cards.NewEmptyCard()
	} else if IsFoundation(from) {
		if card.Value == 0 {
			f.foundations[from.X] = cards.NewCard(cards.VALUE_EMPTY, card.Suite)
		} else {
			f.foundations[from.X] = cards.NewCard(card.Value-1, card.Suite)
		}
	} else {
		if from.Y == 0 {
			f.cascades[from.X][from.Y] = cards.NewEmptyCard()
		} else {
			f.cascades[from.X][from.Y] = nil
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
		selected := f.selected != nil && f.selected.Y == -1 && f.selected.X == i
		hovered := f.position.Y == -1 && f.position.X == i
		foundationsView[i] = c.View(selected, hovered, false)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, foundationsView...)
}

func (f *FreeCell) viewFreeCells() string {
	freecellsView := make([]string, len(f.freecells))

	for i, c := range f.freecells {
		selected := f.selected != nil && f.selected.Y == -1 && f.selected.X == i+Suites
		hovered := f.position.Y == -1 && f.position.X == i+Suites
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
			selected := f.selected != nil && f.selected.Y == j && f.selected.X == i
			hovered := f.position.Y == j && f.position.X == i
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
	posX := f.position.X
	posY := f.position.Y

	// Foundations are not selectable
	if IsFoundation(f.position) {
		return false
	}

	card := f.getCard(f.position)

	// Empty freecells are not selectable
	if IsFreeCell(f.position) {
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
