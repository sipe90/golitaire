package cards

import (
	"fmt"
	"strconv"
	"strings"
)

const width = 6
const height = 5

var (
	values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suites = []string{"♣", "♦", "♥", "♠"}
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
	return c.IsRed()
}

func (c Card) String() string {
	return c.Value + c.Suite
}

func (c Card) View(stacked bool) string {
	label := c.Value + c.Suite
	topFormat := "%-" + strconv.Itoa(width-2) + "s"
	topCenter := strings.ReplaceAll(fmt.Sprintf(topFormat, label), " ", "─")
	top := "╭" + topCenter + "╮"

	if stacked {
		return top
	}

	middle := "│" + strings.Repeat(" ", width-2) + "│"
	bottomFormat := "%" + strconv.Itoa(width-2) + "s"
	bottomCenter := strings.ReplaceAll(fmt.Sprintf(bottomFormat, label), " ", "─")
	bottom := "╰" + bottomCenter + "╯"

	return top + "\n" + strings.Repeat(middle+"\n", height-2) + bottom
}
