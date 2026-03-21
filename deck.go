package main

func FullDeck() []Card {
	deck := make([]Card, 0, 52)

	for suit := Spades; suit <= Diamonds; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}

	return deck
}

func RemainingCards(exclude []Card) []Card {
	excluded := make(map[Card]bool, len(exclude))
	for _, c := range exclude {
		excluded[c] = true
	}

	remaining := make([]Card, 0, 52-len(exclude))
	for _, c := range FullDeck() {
		if !excluded[c] {
			remaining = append(remaining, c)
		}
	}

	return remaining
}
