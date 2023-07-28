package cards

type Deck []Card

func CreateDeck() Deck {
	deck := make(Deck, 0, 52)
	for _, v := range values {
		for _, s := range suites {
			deck = append(deck, CreateCard(&v, &s))
		}
	}

	return deck
}
