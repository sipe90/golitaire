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
	base          lipgloss.Style
	black         lipgloss.Style
	blackSelected lipgloss.Style
	red           lipgloss.Style
	redSelected   lipgloss.Style
	selected      lipgloss.Style
}

const (
	width  = 6
	height = 5

	blackColor         = lipgloss.Color("15")
	selectedBlackColor = lipgloss.Color("0")
	redColor           = lipgloss.Color("9")
	selectedRedColor   = lipgloss.Color("9")
	selectedColor      = lipgloss.Color("0")
	selectedBgColor    = lipgloss.Color("7")
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
			base:          lipgloss.NewStyle(),
			black:         lipgloss.NewStyle().Foreground(blackColor),
			blackSelected: lipgloss.NewStyle().Foreground(selectedBlackColor).Background(selectedBgColor),
			red:           lipgloss.NewStyle().Foreground(redColor),
			redSelected:   lipgloss.NewStyle().Foreground(selectedRedColor).Background(selectedBgColor),
			selected:      lipgloss.NewStyle().Foreground(selectedColor).Background(selectedBgColor),
		},
	}
)

type Card struct {
	Value     string
	Suite     string
	IsVisible bool
}

func CreateCard(value, suite *string) Card {
	return Card{
		Value:     *value,
		Suite:     *suite,
		IsVisible: false,
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

func (c Card) View(selected bool, hovered bool, stacked bool) string {
	var cardStyle lipgloss.Style
	var labelStyle lipgloss.Style

	if selected || hovered {
		cardStyle = pallet.Styles.selected
		if c.IsBlack() {
			labelStyle = pallet.Styles.blackSelected
		} else {
			labelStyle = pallet.Styles.redSelected
		}
	} else {
		cardStyle = pallet.Styles.base
		if c.IsBlack() {
			labelStyle = pallet.Styles.black
		} else {
			labelStyle = pallet.Styles.red
		}
	}

	label := c.Value + c.Suite

	var vpb strings.Builder
	for i := 0; i < width-2-utf8.RuneCountInString(label); i++ {
		fmt.Fprint(&vpb, pallet.Parts.v)
	}

	verticalPadding := cardStyle.Render(vpb.String())
	labelView := labelStyle.Render(label)

	var b strings.Builder

	fmt.Fprint(&b, cardStyle.Render(pallet.Parts.tl))
	fmt.Fprint(&b, labelView)
	fmt.Fprint(&b, cardStyle.Render(verticalPadding))
	fmt.Fprint(&b, cardStyle.Render(pallet.Parts.tr))

	if stacked {
		return b.String()
	}

	var mb strings.Builder

	fmt.Fprint(&mb, cardStyle.Render(pallet.Parts.h))
	for i := 0; i < width-2; i++ {
		fmt.Fprint(&mb, cardStyle.Render(pallet.Parts.m))
	}
	fmt.Fprint(&mb, cardStyle.Render(pallet.Parts.h))

	middle := mb.String()

	fmt.Fprintln(&b)
	for i := 0; i < height-2; i++ {
		fmt.Fprint(&b, middle)
		fmt.Fprintln(&b)
	}

	fmt.Fprint(&b, cardStyle.Render(pallet.Parts.bl))
	fmt.Fprint(&b, verticalPadding)
	fmt.Fprint(&b, labelView)
	fmt.Fprint(&b, cardStyle.Render(pallet.Parts.br))

	return b.String()
}
