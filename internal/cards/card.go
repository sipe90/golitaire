package cards

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
