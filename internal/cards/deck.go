package cards

type Deck []*Card

func CreateDeck() Deck {
	deck := make(Deck, 0, len(values)*len(suites))
	for v := 0; v < len(values); v++ {
		for s := 0; s < len(suites); s++ {
			deck = append(deck, NewCard(Value(v), Suite(s)))
		}
	}

	return deck
}
