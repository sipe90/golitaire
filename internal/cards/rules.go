package cards

type StackType uint8

const (
	SAME_COLOR_ASC  StackType = 0
	SAME_COLOR_DESC StackType = 1
	ALT_COLOR_ASC   StackType = 2
	ALT_COLOR_DESC  StackType = 3
	SAME_SUITE_ASC  StackType = 4
	SAME_SUITE_DESC StackType = 5
)

func CanStack(stackType StackType, from *Card, to *Card) bool {
	switch stackType {
	case SAME_COLOR_ASC:
		return isAsc(from, to) && from.IsBlack() == to.IsBlack()
	case SAME_COLOR_DESC:
		return isDesc(from, to) && from.IsBlack() == to.IsBlack()
	case ALT_COLOR_ASC:
		return isAsc(from, to) && from.IsBlack() != to.IsBlack()
	case ALT_COLOR_DESC:
		return isDesc(from, to) && from.IsBlack() != to.IsBlack()
	case SAME_SUITE_ASC:
		return isAsc(from, to) && from.Suite == to.Suite
	case SAME_SUITE_DESC:
		return isDesc(from, to) && from.Suite == to.Suite
	default:
		return false
	}
}

func CanPlaceOnFoundation(card *Card, foundation *Card) bool {
	if card.Suite != foundation.Suite {
		return false
	}
	if foundation.Value == VALUE_EMPTY {
		return card.Value == 0
	}
	return isAsc(card, foundation)
}

func isAsc(card *Card, to *Card) bool {
	return card.Value-1 == to.Value
}

func isDesc(card *Card, to *Card) bool {
	return card.Value+1 == to.Value
}
