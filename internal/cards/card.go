package cards

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

type Pallet struct {
	Parts  parts
	Styles styles
}

type parts struct {
	tl string
	tr string
	bl string
	br string
	h  string
	v  string
	m  string
}

type styles struct {
	black lipgloss.Style
	red   lipgloss.Style
}

const (
	width  = 6
	height = 5

	blackColor = lipgloss.Color("15")
	redColor   = lipgloss.Color("9")
)

var (
	values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suites = []string{"♣", "♦", "♥", "♠"}

	pallet = Pallet{
		Parts: parts{
			tl: "╭",
			tr: "╮",
			bl: "╰",
			br: "╯",
			h:  "│",
			v:  "─",
			m:  " ",
		},
		Styles: styles{
			black: lipgloss.NewStyle().Foreground(blackColor),
			red:   lipgloss.NewStyle().Foreground(redColor),
		},
	}
)

type Card struct {
	Value      string
	Suite      string
	IsVisible  bool
	IsSelected bool
}

func CreateCard(value, suite *string) Card {
	return Card{
		Value:      *value,
		Suite:      *suite,
		IsVisible:  false,
		IsSelected: false,
	}
}

func (c Card) IsEmpty() bool {
	return c.Suite == "" || c.Value == ""
}

func (c Card) IsRed() bool {
	return c.Suite == "♦" || c.Suite == "♥"
}

func (c Card) IsBlack() bool {
	return !c.IsRed()
}

func (c Card) String() string {
	return c.Value + c.Suite
}

func (c Card) View(stacked bool) string {
	var labelStyle lipgloss.Style
	if c.IsBlack() {
		labelStyle = pallet.Styles.black
	} else {
		labelStyle = pallet.Styles.red
	}

	label := c.Value + c.Suite

	var hpb strings.Builder
	for i := 0; i < width-2-utf8.RuneCountInString(label); i++ {
		fmt.Fprint(&hpb, pallet.Parts.v)
	}

	verticalPadding := hpb.String()
	labelView := labelStyle.Render(label)

	var b strings.Builder

	fmt.Fprint(&b, pallet.Parts.tl)
	fmt.Fprint(&b, labelView)
	fmt.Fprint(&b, verticalPadding)
	fmt.Fprint(&b, pallet.Parts.tr)

	if stacked {
		return b.String()
	}

	var mb strings.Builder

	fmt.Fprint(&mb, pallet.Parts.h)
	for i := 0; i < width-2; i++ {
		fmt.Fprint(&mb, pallet.Parts.m)
	}
	fmt.Fprint(&mb, pallet.Parts.h)

	middle := mb.String()

	fmt.Fprintln(&b)
	for i := 0; i < height-2; i++ {
		fmt.Fprint(&b, middle)
		fmt.Fprintln(&b)
	}

	fmt.Fprint(&b, pallet.Parts.bl)
	fmt.Fprint(&b, verticalPadding)
	fmt.Fprint(&b, label)
	fmt.Fprint(&b, pallet.Parts.br)

	return b.String()
}
