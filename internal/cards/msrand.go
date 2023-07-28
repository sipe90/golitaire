package cards

import "math"

const rMax32 = math.MaxInt32

var seed int

func rnd() int {
	seed = (seed*214013 + 2531011) & rMax32
	return seed >> 16
}

func Shuffle(d *Deck, s int) Deck {
	seed = s
	temp := make(Deck, len(*d))
	shuffled := make(Deck, 0, len(*d))
	copy(temp, *d)

	for i := len(temp); i > 0; i-- {
		index := rnd() % i
		shuffled = append(shuffled, temp[index])
		temp[index] = temp[i-1]
	}

	return shuffled
}
